package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/pointer"
	"github.com/ekubyshin/metrics_agent/internal/storage"
)

const (
	contentTypeHeader = "Content-Type"
	applicationJSON   = "application/json"
)

type RestHandler struct {
	route string
	db    storage.Storage[metrics.MetricsKey, metrics.Metrics]
}

func NewRestHandler(
	db storage.Storage[metrics.MetricsKey, metrics.Metrics]) *RestHandler {
	return &RestHandler{
		route: "/",
		db:    db,
	}
}

func (h *RestHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(contentTypeHeader, applicationJSON)
	if !checkContentType(r) {
		w.WriteHeader(http.StatusBadRequest)
	}
	ms, ok := parseMetircs(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if ms.MType == handlers.GaugeActionKey {
		if ok := h.putGauge(r.Context(), ms); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		var newCounterVal metrics.Counter
		if newCounterVal, ok = h.putCounter(r.Context(), *ms); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			ms.Delta = pointer.From[int64](int64(newCounterVal))
		}
	}
	res, err := json.MarshalIndent(ms, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func (h *RestHandler) Value(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(contentTypeHeader, applicationJSON)
	if !checkContentType(r) {
		w.WriteHeader(http.StatusBadRequest)
	}
	metrics, ok := parseMetircs(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if metrics.MType == handlers.GaugeActionKey {
		val, ok := h.db.Get(r.Context(), metrics.Key())
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metrics.Value = pointer.From[float64](float64(*val.Value))
	} else {
		val, ok := h.db.Get(r.Context(), metrics.Key())
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metrics.Delta = pointer.From[int64](int64(*val.Delta))
	}
	res, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func (h *RestHandler) putGauge(ctx context.Context, m *metrics.Metrics) bool {
	if m.Value == nil {
		return false
	}
	err := h.db.Put(ctx, m.Key(), *m)
	return err == nil
}

func (h *RestHandler) putCounter(ctx context.Context, m metrics.Metrics) (metrics.Counter, bool) {
	if m.Delta == nil {
		return 0, false
	}
	nv := *m.Delta
	if v, ok := h.db.Get(ctx, m.Key()); ok {
		prev := int64(0)
		if v.Delta != nil {
			prev = *v.Delta
		}
		nv += prev
		m.Delta = pointer.From[int64](int64(nv))
	}
	err := h.db.Put(ctx, m.Key(), m)
	if err != nil {
		return 0, false
	}
	return metrics.Counter(nv), true
}

func checkContentType(r *http.Request) bool {
	return r.Header.Get(contentTypeHeader) == applicationJSON
}

func parseMetircs(r *http.Request) (*metrics.Metrics, bool) {
	var buf bytes.Buffer
	var metrics metrics.Metrics
	defer r.Body.Close()
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return nil, false
	}
	if err := json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		return nil, false
	}
	if metrics.MType != handlers.CounterActionKey && metrics.MType != handlers.GaugeActionKey {
		return &metrics, false
	}
	return &metrics, true
}

func (h *RestHandler) BaseURL() string {
	return h.route
}
