package server

import (
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/handlers/counter"
	"github.com/ekubyshin/metrics_agent/internal/handlers/gauge"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Server interface {
	Run() error
}

type ChiServer struct {
	router *chi.Mux
}

func NewServer() Server {
	db := storage.NewMemoryStorage()
	router := chi.NewRouter()
	gaugePostHandler := gauge.NewGaugePostHandler(db)
	counterPostHandler := counter.NewCounterPostHandler(db)
	gaugeGetHandler := gauge.NewGaugeGetHandler(db)
	counterGetHanlder := counter.NewCounterGetHandler(db)
	router.Route("/update", func(r chi.Router) {
		r.Post(gaugePostHandler.BaseURL(), gaugePostHandler.ServeHTTP)
		r.Post(counterPostHandler.BaseURL(), counterPostHandler.ServeHTTP)
	})
	router.Route("/value", func(r chi.Router) {
		r.Get(gaugeGetHandler.BaseURL(), gaugeGetHandler.ServeHTTP)
		r.Get(counterGetHanlder.BaseURL(), counterGetHanlder.ServeHTTP)
	})
	return &ChiServer{
		router: router,
	}
}

func (s *ChiServer) Run() error {
	return http.ListenAndServe(":8080", s.router)
}
