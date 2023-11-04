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

type CounterPostHandler struct {
	route string
	db    storage.Storage[metrics.MetricsKey, metrics.Metrics]
}

func NewCounterPostHandler(db storage.Storage[metrics.MetricsKey, metrics.Metrics]) handlers.Handler {
	return &CounterPostHandler{
		route: "/counter/{name}/{value}",
		db:    db,
	}
}

func (m *CounterPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (m *CounterPostHandler) BaseURL() string {
	return m.route
}
