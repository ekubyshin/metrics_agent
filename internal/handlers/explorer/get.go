package explorer

import (
	"encoding/json"
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/storage"
)

type ExplorerHandler struct {
	route string
	db    storage.Storage[handlers.Key, any]
}

func NewExplorerHandler(db storage.Storage[handlers.Key, any]) handlers.Handler {
	return &ExplorerHandler{
		route: "/",
		db:    db,
	}
}

func (e *ExplorerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	elems := e.db.List()
	out := make(map[string]any)
	if len(elems) > 0 {
		for _, v := range elems {
			out[v.Key.Name] = v.Value
		}
	}
	res, err := json.Marshal(out)
	w.Header().Add("Content-type", "application/json")
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
