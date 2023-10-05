package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/ekubyshin/metrics_agent/internal/server/types"
	"github.com/ekubyshin/metrics_agent/internal/utils"
)

type Metrics struct {
	route string
}

type metricsHanlerPath struct {
	metricsType  string
	metricsName  string
	metricsValue string
}

func NewMetricsHandler() *Metrics {
	return &Metrics{
		route: "/update/",
	}
}

func (m *Metrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_, err := m.parsePath(r.URL)
	if err != nil {
		if errors.Is(err, types.NewMetricsHandlerInvalidNameError()) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (m *Metrics) Route() string {
	return m.route
}

func (m *Metrics) parsePath(url *url.URL) (*metricsHanlerPath, error) {
	parts := strings.Split(url.Path, "/")
	parts = utils.DeleteEmpty(parts)
	if len(parts) > 4 {
		return nil, types.NewUnknowSchemeError()
	}
	switch len(parts) {
	case 1:
		return nil, types.NewMetricsHandlerInvalidTypeError()
	case 2:
		return nil, types.NewMetricsHandlerInvalidNameError()
	case 3:
		return nil, types.NewInvalidMetricsValue()
	}

	return &metricsHanlerPath{
		metricsType:  parts[1],
		metricsName:  parts[2],
		metricsValue: parts[3],
	}, nil
}
