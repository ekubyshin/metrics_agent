package gauge

import (
	"net/http"
	"strconv"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/go-chi/chi/v5"
)

const (
	GetURL  = "/gauge/{name}"
	PostURL = "/gauge/{name}/{value}"
)

type GaugeHandler struct {
	db storage.Storage[metrics.MetricsKey, metrics.Metrics]
}

func NewGaugeHandler(db storage.Storage[metrics.MetricsKey, metrics.Metrics]) *GaugeHandler {
	return &GaugeHandler{
		db: db,
	}
}

func (m *GaugeHandler) Get(w http.ResponseWriter, r *http.Request) {
	paramName := chi.URLParam(r, handlers.ParamNameKey)
	if paramName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if v, ok := m.db.Get(
		r.Context(),
		metrics.MetricsKey{
			ID:   paramName,
			Type: handlers.GaugeActionKey,
		}); ok {
		_, err := w.Write([]byte(strconv.FormatFloat(float64(*v.Value), 'f', -1, 64)))
		if err == nil {
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (m *GaugeHandler) Post(w http.ResponseWriter, r *http.Request) {
	paramName := chi.URLParam(r, handlers.ParamNameKey)
	paramValue := chi.URLParam(r, handlers.ParamValueKey)
	parsedValue, err := strconv.ParseFloat(paramValue, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = m.db.Put(
		r.Context(),
		metrics.MetricsKey{
			ID:   paramName,
			Type: handlers.GaugeActionKey,
		},
		metrics.Metrics{
			ID:    paramName,
			MType: handlers.GaugeActionKey,
			Value: &parsedValue,
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
