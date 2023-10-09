package counter

import (
	"net/http"
	"strconv"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/go-chi/chi/v5"
)

type CounterPostHandler struct {
	route string
	db    storage.Storage[string, types.Counter]
}

func NewCounterPostHandler(db storage.Storage[string, types.Counter]) handlers.Handler {
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
	if v, ok := m.db.Get(paramName); ok {
		m.db.Put(paramName, types.Counter(parsedValue)+v)
	} else {
		m.db.Put(paramName, types.Counter(parsedValue))
	}
	w.WriteHeader(http.StatusOK)
}

func (m *CounterPostHandler) BaseURL() string {
	return m.route
}
