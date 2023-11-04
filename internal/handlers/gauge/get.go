package gauge

import (
	"net/http"
	"strconv"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/go-chi/chi/v5"
)

type GaugeGetHandler struct {
	route string
	db    storage.Storage[metrics.MetricsKey, metrics.Metrics]
}

func NewGaugeGetHandler(db storage.Storage[metrics.MetricsKey, metrics.Metrics]) *GaugeGetHandler {
	return &GaugeGetHandler{
		route: "/gauge/{name}",
		db:    db,
	}
}

func (m *GaugeGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paramName := chi.URLParam(r, handlers.ParamNameKey)
	if paramName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if v, ok := m.db.Get(
		metrics.MetricsKey{
			ID:    paramName,
			MType: handlers.GaugeActionKey,
		}); ok {
		_, err := w.Write([]byte(strconv.FormatFloat(float64(*v.Value), 'f', -1, 64)))
		if err == nil {
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (m *GaugeGetHandler) BaseURL() string {
	return m.route
}
