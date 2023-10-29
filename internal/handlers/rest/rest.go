package rest

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/ekubyshin/metrics_agent/internal/utils"
)

const (
	contentTypeHeader = "Content-Type"
	applicationJSON   = "application/json"
)

type RestHandler struct {
	route     string
	dbCounter storage.Storage[string, types.Counter]
	dbGauge   storage.Storage[string, types.Gauge]
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewRestHandler(
	dbCounter storage.Storage[string, types.Counter],
	dbGauge storage.Storage[string, types.Gauge]) *RestHandler {
	return &RestHandler{
		route:     "/",
		dbCounter: dbCounter,
		dbGauge:   dbGauge,
	}
}

func (h *RestHandler) Update(w http.ResponseWriter, r *http.Request) {
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
		if ok := h.putGauge(metrics); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		var newCounterVal types.Counter
		if newCounterVal, ok = h.putCounter(metrics); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			metrics.Delta = utils.ToPointer[int64](int64(newCounterVal))
		}
	}
	res, err := json.MarshalIndent(metrics, "", "  ")
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
		val, ok := h.dbGauge.Get(metrics.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metrics.Value = utils.ToPointer[float64](float64(val))
	} else {
		val, ok := h.dbCounter.Get(metrics.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metrics.Delta = utils.ToPointer[int64](int64(val))
	}
	res, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func (h *RestHandler) putGauge(m *Metrics) bool {
	if m.Value == nil {
		return false
	}
	h.dbGauge.Put(m.ID, types.Gauge(*m.Value))
	return true
}

func (h *RestHandler) putCounter(m *Metrics) (types.Counter, bool) {
	if m.Delta == nil {
		return 0, false
	}
	nv := types.Counter(*m.Delta)
	if v, ok := h.dbCounter.Get(m.ID); ok {
		nv += v
	}
	h.dbCounter.Put(m.ID, nv)
	return nv, true
}

func checkContentType(r *http.Request) bool {
	return r.Header.Get(contentTypeHeader) == applicationJSON
}

func parseMetircs(r *http.Request) (*Metrics, bool) {
	var buf bytes.Buffer
	var metrics Metrics
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
