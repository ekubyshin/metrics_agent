package explorer

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/storage"
)

type ExplorerHandler struct {
	db storage.Storage[metrics.MetricsKey, metrics.Metrics]
}

func NewExplorerHandler(db storage.Storage[metrics.MetricsKey, metrics.Metrics]) *ExplorerHandler {
	return &ExplorerHandler{
		db: db,
	}
}

func (e *ExplorerHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	elems, err := e.db.List(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(elems) == 0 {
		_, _ = w.Write([]byte(`{}`))
		return
	}
	out := make(map[string]any)
	for _, v := range elems {
		if v.Value.Value == nil {
			out[v.Key.Type+"_"+v.Key.ID] = v.Value.Delta
			continue
		}
		out[v.Key.Type+"_"+v.Key.ID] = v.Value.Value
	}
	res, err := json.Marshal(out)
	if err == nil {
		_, err = w.Write(res)
		if err == nil {
			return
		}
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}
