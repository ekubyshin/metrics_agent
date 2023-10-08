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
	db    storage.Storage[handlers.Key, any]
}

func NewCounterPostHandler(db storage.Storage[handlers.Key, any]) handlers.Handler {
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
	key := handlers.Key{Type: handlers.CounterActionKey, Name: paramName}
	if v, ok := m.db.Get(key); ok {
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
