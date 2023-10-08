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
	db    storage.Storage[handlers.Key, any]
}

func NewGaugeGetHandler(db storage.Storage[handlers.Key, any]) handlers.Handler {
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
	if rv, ok := m.db.Get(handlers.Key{Type: handlers.GaugeActionKey, Name: paramName}); ok {
		if v, ok2 := rv.(types.Gauge); ok2 {
			_, err := w.Write([]byte(strconv.FormatFloat(float64(v), 'f', -1, 64)))
			if err == nil {
				return
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (m *GaugeGetHandler) BaseURL() string {
	return m.route
}
