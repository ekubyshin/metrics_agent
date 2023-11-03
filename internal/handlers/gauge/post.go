package gauge

import (
	"net/http"
	"strconv"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/go-chi/chi/v5"
)

type GaugePostHandler struct {
	route string
	db    storage.Storage[types.MetricsKey, types.Metrics]
}

func NewGaugePostHandler(db storage.Storage[types.MetricsKey, types.Metrics]) *GaugePostHandler {
	return &GaugePostHandler{
		route: "/gauge/{name}/{value}",
		db:    db,
	}
}

func (m *GaugePostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paramName := chi.URLParam(r, handlers.ParamNameKey)
	paramValue := chi.URLParam(r, handlers.ParamValueKey)
	parsedValue, err := strconv.ParseFloat(paramValue, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m.db.Put(
		types.MetricsKey{
			ID:    paramName,
			MType: handlers.GaugeActionKey,
		},
		types.Metrics{
			ID:    paramName,
			MType: handlers.GaugeActionKey,
			Value: &parsedValue,
		},
	)
	w.WriteHeader(http.StatusOK)
}

func (m *GaugePostHandler) BaseURL() string {
	return m.route
}
