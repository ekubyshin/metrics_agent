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
	db    storage.Storage
}

func NewCounterPostHandler(db storage.Storage) handlers.Handler {
	return &CounterPostHandler{
		route: "/counter/{name}/{value:^[0-9]+}",
		db:    db,
	}
}

func (m *CounterPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paramName := chi.URLParam(r, "name")
	paramValue := chi.URLParam(r, "value")
	if paramName == "" || paramValue == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	parsedValue, err := strconv.ParseInt(paramValue, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	key := handlers.Key{Type: "counter", Name: paramName}
	if v, ok := m.db.Get(key); ok == nil {
		if pv, ok2 := v.(types.Counter); ok2 {
			m.db.Put(key, types.Counter(parsedValue)+pv)
		}
	} else {
		m.db.Put(key, types.Counter(parsedValue))
	}
	w.WriteHeader(http.StatusOK)
}

func (m *CounterPostHandler) BaseURL() string {
	return m.route
}
