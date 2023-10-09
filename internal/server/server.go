package server

import (
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/handlers/counter"
	"github.com/ekubyshin/metrics_agent/internal/handlers/explorer"
	"github.com/ekubyshin/metrics_agent/internal/handlers/gauge"
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

func NewServer(endpoint config.Address) Server {
	dbCounter := storage.NewMemoryStorage[string, types.Counter]()
	dbGauge := storage.NewMemoryStorage[string, types.Gauge]()
	router := chi.NewRouter()
	gaugePostHandler := gauge.NewGaugePostHandler(dbGauge)
	counterPostHandler := counter.NewCounterPostHandler(dbCounter)
	gaugeGetHandler := gauge.NewGaugeGetHandler(dbGauge)
	counterGetHanlder := counter.NewCounterGetHandler(dbCounter)
	listHanlder := explorer.NewExplorerHandler(dbCounter, dbGauge)
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
	return &ChiServer{
		router:   router,
		endpoint: endpoint,
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
