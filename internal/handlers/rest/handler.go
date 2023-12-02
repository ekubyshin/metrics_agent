package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/pointer"
	"github.com/ekubyshin/metrics_agent/internal/storage"
)

type RestHandler struct {
	db storage.Storage[metrics.MetricsKey, metrics.Metrics]
}

func NewRestHandler(
	db storage.Storage[metrics.MetricsKey, metrics.Metrics]) *RestHandler {
	return &RestHandler{
		db: db,
	}
}

type batchMetrics struct {
	Metrics []metrics.Metrics `json:"metrics"`
}

func (h *RestHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(handlers.ContentTypeHeader, handlers.ApplicationJSON)
	if !checkContentType(r) {
		w.WriteHeader(http.StatusBadRequest)
	}
	ms, ok := parseSingleMetric(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	nm, ok := h.putMetric(ctx, ms)
	if !ok || nm == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.MarshalIndent(nm, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(res)
}

func (h *RestHandler) Updates(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(handlers.ContentTypeHeader, handlers.ApplicationJSON)
	if !checkContentType(r) {
		w.WriteHeader(http.StatusBadRequest)
	}
	ms, ok := parseMetrics(r)
	if !ok || len(ms) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	kv := make([]storage.KeyValuer[metrics.MetricsKey, metrics.Metrics], 0, len(ms))
	for _, m := range ms {
		kv = append(kv, storage.KeyValuer[metrics.MetricsKey, metrics.Metrics]{Key: m.Key(), Value: m})
	}
	out, err := h.db.PutBatch(ctx, kv)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res := batchMetrics{Metrics: out}
	s, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(s)
}

func (h *RestHandler) Value(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(handlers.ContentTypeHeader, handlers.ApplicationJSON)
	if !checkContentType(r) {
		w.WriteHeader(http.StatusBadRequest)
	}
	metrics, ok := parseSingleMetric(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if metrics.MType == handlers.GaugeActionKey {
		val, ok := h.db.Get(ctx, metrics.Key())
		if !ok || val.Value == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metrics.Value = pointer.From[float64](float64(*val.Value))
	} else {
		val, ok := h.db.Get(ctx, metrics.Key())
		if !ok || val.Delta == nil {
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
	_, _ = w.Write(res)
}

func (h *RestHandler) putGauge(ctx context.Context, m *metrics.Metrics) bool {
	if m.Value == nil {
		return false
	}
	_, err := h.db.Put(ctx, m.Key(), *m)
	return err == nil
}

func (h *RestHandler) putCounter(ctx context.Context, m *metrics.Metrics) (metrics.Counter, bool) {
	if m == nil || m.Delta == nil {
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
	_, err := h.db.Put(ctx, m.Key(), *m)
	if err != nil {
		return 0, false
	}
	return metrics.Counter(nv), true
}

func (h *RestHandler) putMetric(ctx context.Context, m *metrics.Metrics) (*metrics.Metrics, bool) {
	if m == nil {
		return nil, false
	}
	if m.MType == handlers.GaugeActionKey {
		if ok := h.putGauge(ctx, m); !ok {
			return m, false
		}
		return m, true
	}
	var newCounterVal metrics.Counter
	var ok bool
	if newCounterVal, ok = h.putCounter(ctx, m); !ok {
		return m, false
	}
	m.Delta = pointer.From[int64](int64(newCounterVal))
	return m, true
}

func checkContentType(r *http.Request) bool {
	return r.Header.Get(handlers.ContentTypeHeader) == handlers.ApplicationJSON
}

func parseSingleMetric(r *http.Request) (*metrics.Metrics, bool) {
	var buf bytes.Buffer
	metrics := &metrics.Metrics{}
	defer r.Body.Close()
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return nil, false
	}
	if err := json.Unmarshal(buf.Bytes(), metrics); err != nil {
		return nil, false
	}
	if metrics.MType != handlers.CounterActionKey && metrics.MType != handlers.GaugeActionKey {
		return metrics, false
	}
	return metrics, true
}

func parseMetrics(r *http.Request) ([]metrics.Metrics, bool) {
	var buf bytes.Buffer
	var elems []metrics.Metrics
	defer r.Body.Close()
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return nil, false
	}
	if err := json.Unmarshal(buf.Bytes(), &elems); err != nil {
		return nil, false
	}
	return elems, true
}
