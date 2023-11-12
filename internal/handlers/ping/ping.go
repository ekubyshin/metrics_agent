package ping

import (
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/storage"
)

type PingHandler struct {
	st    storage.Storage[metrics.MetricsKey, metrics.Metrics]
	route string
}

func NewPingHandler(st storage.Storage[metrics.MetricsKey, metrics.Metrics]) *PingHandler {
	return &PingHandler{
		st:    st,
		route: "/ping",
	}
}

func (m *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.st == nil || m.st.Ping() != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *PingHandler) BaseURL() string {
	return m.route
}
