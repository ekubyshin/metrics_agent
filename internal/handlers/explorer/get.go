package explorer

import (
	"encoding/json"
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/storage"
)

type ExplorerHandler struct {
	route string
	db    storage.Storage[metrics.MetricsKey, metrics.Metrics]
}

func NewExplorerHandler(db storage.Storage[metrics.MetricsKey, metrics.Metrics]) handlers.Handler {
	return &ExplorerHandler{
		route: "/",
		db:    db,
	}
}

func (e *ExplorerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")
	elems, err := e.db.List(r.Context())
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

func (e *ExplorerHandler) BaseURL() string {
	return e.route
}
