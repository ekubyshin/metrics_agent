package server

import (
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/ekubyshin/metrics_agent/internal/handlers"
	l "github.com/ekubyshin/metrics_agent/internal/logger"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	mw "github.com/ekubyshin/metrics_agent/internal/middlewares"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Server interface {
	Run() error
}

type ChiServer struct {
	router   *chi.Mux
	endpoint config.Address
}

func NewServer(
	cfg config.Config,
	logger l.Logger,
	st storage.Storage[metrics.MetricsKey, metrics.Metrics],
) *ChiServer {
	router := chi.NewRouter()
	router.Use(mw.NewRequestLogger(logger))
	router.Use(mw.NewResponseLogger(logger))
	router.Use(mw.GzipReader)
	router.Use(mw.GzipHandler)
	RegisterRoutes(router, st)

	return &ChiServer{
		router:   router,
		endpoint: cfg.Address,
	}
}

func (s *ChiServer) Run() error {
	return http.ListenAndServe(s.endpoint.ToString(), s.router)
}

func GetErrorStatusCode(r *http.Request) int {
	t := chi.URLParam(r, handlers.ParamTypeKey)
	switch t {
	case handlers.GaugeActionKey:
		return http.StatusNotFound
	case handlers.CounterActionKey:
		return http.StatusNotFound
	default:
		return http.StatusNotImplemented
	}
}
