package ping

import (
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/storage"
)

type PingHandler struct {
	db    *storage.DBStorage
	route string
}

func NewPingHandler(db *storage.DBStorage) *PingHandler {
	return &PingHandler{
		db:    db,
		route: "/ping",
	}
}

func (m *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.db == nil || m.db.Ping() != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *PingHandler) BaseURL() string {
	return m.route
}
