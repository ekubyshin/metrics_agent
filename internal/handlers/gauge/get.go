package gauge

import (
	"net/http"
	"strconv"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/go-chi/chi/v5"
)

type GaugeGetHandler struct {
	route string
	db    storage.Storage[string, types.Gauge]
}

func NewGaugeGetHandler(db storage.Storage[string, types.Gauge]) *GaugeGetHandler {
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
	if v, ok := m.db.Get(paramName); ok {
		_, err := w.Write([]byte(strconv.FormatFloat(float64(v), 'f', -1, 64)))
		if err == nil {
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (m *GaugeGetHandler) BaseURL() string {
	return m.route
}
