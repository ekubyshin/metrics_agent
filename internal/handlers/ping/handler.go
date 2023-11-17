package ping

import (
	"context"
	"net/http"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/storage"
)

type PingHandler struct {
	st storage.Storage[metrics.MetricsKey, metrics.Metrics]
}

func NewPingHandler(st storage.Storage[metrics.MetricsKey, metrics.Metrics]) *PingHandler {
	return &PingHandler{
		st: st,
	}
}

func (m *PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if m.st == nil || m.st.Ping(ctx) != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *PingHandler) Features(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(handlers.ContentTypeHeader, handlers.ApplicationJSON)
	_, _ = w.Write([]byte(`{"batch":true}`))
}
