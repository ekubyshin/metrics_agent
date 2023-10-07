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
	db    storage.Storage
}

func NewCounterGetHandler(db storage.Storage) handlers.Handler {
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
	if rv, ok := m.db.Get(handlers.Key{Type: "counter", Name: paramName}); ok == nil {
		if v, ok2 := rv.(types.Counter); ok2 {
			c := int64(v)
			_, err := w.Write([]byte(strconv.FormatInt(c, 10)))
			if err == nil {
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("0"))
}

func (m *CounterGetHandler) BaseURL() string {
	return m.route
}
