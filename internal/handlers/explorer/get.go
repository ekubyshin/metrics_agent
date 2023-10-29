package explorer

import (
	"encoding/json"
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/ekubyshin/metrics_agent/internal/types"
)

type ExplorerHandler struct {
	route     string
	dbCounter storage.Storage[string, types.Counter]
	dbGauge   storage.Storage[string, types.Gauge]
}

func NewExplorerHandler(dbCounter storage.Storage[string, types.Counter], dbGauge storage.Storage[string, types.Gauge]) handlers.Handler {
	return &ExplorerHandler{
		route:     "/",
		dbCounter: dbCounter,
		dbGauge:   dbGauge,
	}
}

func (e *ExplorerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	elemsGauge := e.dbGauge.List()
	elemsCounter := e.dbCounter.List()
	out := make(map[string]any)
	for _, v := range elemsGauge {
		out[handlers.GaugeActionKey+"_"+v.Key] = v.Value
	}
	for _, v := range elemsCounter {
		out[handlers.CounterActionKey+"_"+v.Key] = v.Value
	}
	res, err := json.Marshal(out)
	w.Header().Add("Content-type", "text/html")
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
