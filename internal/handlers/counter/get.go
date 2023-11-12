package counter

import (
	"net/http"
	"strconv"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/pointer"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/go-chi/chi/v5"
)

const (
	GetURL  = "/counter/{name}"
	PostURL = "/counter/{name}/{value}"
)

type CounterHandler struct {
	db storage.Storage[metrics.MetricsKey, metrics.Metrics]
}

func NewCounterHandler(db storage.Storage[metrics.MetricsKey, metrics.Metrics]) *CounterHandler {
	return &CounterHandler{
		db: db,
	}
}

func (m *CounterHandler) Get(w http.ResponseWriter, r *http.Request) {
	paramName := chi.URLParam(r, "name")
	if paramName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if v, ok := m.db.Get(metrics.MetricsKey{ID: paramName, MType: handlers.CounterActionKey}); ok {
		_, err := w.Write([]byte(strconv.FormatInt(*v.Delta, 10)))
		if err == nil {
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (m *CounterHandler) Post(w http.ResponseWriter, r *http.Request) {
	paramName := chi.URLParam(r, handlers.ParamNameKey)
	paramValue := chi.URLParam(r, handlers.ParamValueKey)
	if paramName == "" || paramValue == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	parsedValue, err := strconv.ParseInt(paramValue, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	key := metrics.MetricsKey{ID: paramName, MType: handlers.CounterActionKey}
	if v, ok := m.db.Get(key); ok {
		prev := int64(0)
		if v.Delta != nil {
			prev = *v.Delta
		}
		v.Delta = pointer.From[int64](prev + parsedValue)
		m.db.Put(key, v)
	} else {
		m.db.Put(key, metrics.Metrics{
			ID:    paramName,
			MType: handlers.CounterActionKey,
			Delta: pointer.From[int64](parsedValue),
		})
	}
	w.WriteHeader(http.StatusOK)
}
