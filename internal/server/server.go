package server

import (
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/handlers/counter"
	"github.com/ekubyshin/metrics_agent/internal/handlers/explorer"
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
	listHanlder := explorer.NewExplorerHandler(db)
	router.Get(listHanlder.BaseURL(), listHanlder.ServeHTTP)
	router.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		t := chi.URLParam(r, "type")
		switch t {
		case "gauge":
			gaugePostHandler.ServeHTTP(w, r)
		case "counter":
			counterPostHandler.ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	})
	router.Post("/update/{type}/", func(w http.ResponseWriter, r *http.Request) {
		t := chi.URLParam(r, "type")
		switch t {
		case "gauge":
			w.WriteHeader(http.StatusNotFound)
		case "counter":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	})
	// router.Route("/update", func(r chi.Router) {
	// 	r.Post("/{a}/{b}/", func(w http.ResponseWriter, r *http.Request) {
	// 		b := chi.URLParam(r, "b")
	// 		if b != "gauge" && b != "counter" {
	// 			w.WriteHeader(http.StatusNotImplemented)
	// 			return
	// 		}
	// 		w.WriteHeader(http.StatusNotFound)
	// 	})
	// 	r.Post(gaugePostHandler.BaseURL(), gaugePostHandler.ServeHTTP)
	// 	r.Post(counterPostHandler.BaseURL(), counterPostHandler.ServeHTTP)
	// })
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
