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
	db    storage.Storage
}

func NewGaugePostHandler(db storage.Storage) handlers.Handler {
	return &GaugePostHandler{
		route: "/gauge/{name}/{value}",
		db:    db,
	}
}

func (m *GaugePostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paramName := chi.URLParam(r, "name")
	paramValue := chi.URLParam(r, "value")
	parsedValue, err := strconv.ParseFloat(paramValue, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m.db.Put(handlers.Key{Type: "gauge", Name: paramName}, types.Gauge(parsedValue))
	w.WriteHeader(http.StatusOK)
}

func (m *GaugePostHandler) BaseURL() string {
	return m.route
}
