package server

import (
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/go-chi/chi/v5"
)

type Server interface {
	Run() error
}

type ChiServer struct {
	router *chi.Mux
}

func NewServer() Server {
	router := chi.NewRouter()
	gaugeHandler := handlers.NewGaugeHandler()
	counterHandler := handlers.NewCounterHandler()
	router.Route("/update", func(r chi.Router) {
		r.Post(gaugeHandler.BaseURL(), gaugeHandler.ServeHTTP)
		r.Post(counterHandler.BaseURL(), counterHandler.ServeHTTP)
	})
	return &ChiServer{
		router: router,
	}
}

func (s *ChiServer) Run() error {
	return http.ListenAndServe(":8080", s.router)
}
