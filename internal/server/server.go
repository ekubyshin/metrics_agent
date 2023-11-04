package server

import (
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/handlers/counter"
	"github.com/ekubyshin/metrics_agent/internal/handlers/explorer"
	"github.com/ekubyshin/metrics_agent/internal/handlers/gauge"
	"github.com/ekubyshin/metrics_agent/internal/handlers/rest"
	l "github.com/ekubyshin/metrics_agent/internal/logger"
	mw "github.com/ekubyshin/metrics_agent/internal/middlewares"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/go-chi/chi/v5"
)

type Server interface {
	Run() error
}

type ChiServer struct {
	router   *chi.Mux
	endpoint config.Address
}

func NewServer(cfg config.Config, logger l.Logger) *ChiServer {
	db := storage.NewMemoryStorage[types.MetricsKey, types.Metrics]()
	var w *storage.FileStorage[types.MetricsKey, types.Metrics]
	var err error
	if cfg.FileStoragePath != nil && *cfg.FileStoragePath != "" {
		w, err = storage.NewFileStorage(db, *cfg.FileStoragePath, *cfg.Restore, cfg.StoreDuration())
		if err != nil {
			panic(err)
		}
	}
	router := chi.NewRouter()
	router.Use(mw.NewRequestLogger(logger))
	router.Use(mw.NewResponseLogger(logger))
	router.Use(mw.GzipReader)
	router.Use(mw.GzipHandler)
	if w != nil {
		registerRoutes(router, w)
	} else {
		registerRoutes(router, db)
	}

	return &ChiServer{
		router:   router,
		endpoint: cfg.Address,
	}
}

func registerRoutes(
	router *chi.Mux,
	db storage.Storage[types.MetricsKey, types.Metrics]) {
	gaugePostHandler := gauge.NewGaugePostHandler(db)
	counterPostHandler := counter.NewCounterPostHandler(db)
	gaugeGetHandler := gauge.NewGaugeGetHandler(db)
	counterGetHanlder := counter.NewCounterGetHandler(db)
	listHanlder := explorer.NewExplorerHandler(db)
	restHandler := rest.NewRestHandler(db)
	router.Get(listHanlder.BaseURL(), listHanlder.ServeHTTP)
	router.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		t := chi.URLParam(r, handlers.ParamTypeKey)
		switch t {
		case handlers.GaugeActionKey:
			gaugePostHandler.ServeHTTP(w, r)
		case handlers.CounterActionKey:
			counterPostHandler.ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	})
	router.Post("/update/{type}/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(GetErrorStatusCode(r))
	})
	router.Route("/value", func(r chi.Router) {
		r.Get(gaugeGetHandler.BaseURL(), gaugeGetHandler.ServeHTTP)
		r.Get(counterGetHanlder.BaseURL(), counterGetHanlder.ServeHTTP)
	})
	router.Post("/update/", restHandler.Update)
	router.Post("/value/", restHandler.Value)
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
