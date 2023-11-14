package counter

import (
	"context"
	"net/http"
	"strconv"
	"time"

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
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if v, ok := m.db.Get(ctx, metrics.MetricsKey{ID: paramName, Type: handlers.CounterActionKey}); ok {
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
	key := metrics.MetricsKey{ID: paramName, Type: handlers.CounterActionKey}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if v, ok := m.db.Get(ctx, key); ok {
		prev := int64(0)
		if v.Delta != nil {
			prev = *v.Delta
		}
		v.Delta = pointer.From[int64](prev + parsedValue)
		_, err := m.db.Put(ctx, key, v)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		_, err := m.db.Put(
			ctx,
			key,
			metrics.Metrics{
				ID:    paramName,
				MType: handlers.CounterActionKey,
				Delta: pointer.From[int64](parsedValue),
			},
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
