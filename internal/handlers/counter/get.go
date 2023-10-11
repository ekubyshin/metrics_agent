package counter

import (
	"net/http"
	"strconv"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/go-chi/chi/v5"
)

type CounterGetHandler struct {
	route string
	db    storage.Storage[string, types.Counter]
}

func NewCounterGetHandler(db storage.Storage[string, types.Counter]) handlers.Handler {
	return &CounterGetHandler{
		route: "/counter/{name}",
		db:    db,
	}
}

func (m *CounterGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paramName := chi.URLParam(r, "name")
	if paramName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if v, ok := m.db.Get(paramName); ok {
		c := int64(v)
		_, err := w.Write([]byte(strconv.FormatInt(c, 10)))
		if err == nil {
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (m *CounterGetHandler) BaseURL() string {
	return m.route
}
