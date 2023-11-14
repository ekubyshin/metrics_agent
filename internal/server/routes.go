package server

import (
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/handlers/counter"
	"github.com/ekubyshin/metrics_agent/internal/handlers/explorer"
	"github.com/ekubyshin/metrics_agent/internal/handlers/gauge"
	"github.com/ekubyshin/metrics_agent/internal/handlers/ping"
	"github.com/ekubyshin/metrics_agent/internal/handlers/rest"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(
	router *chi.Mux,
	st storage.Storage[metrics.MetricsKey, metrics.Metrics]) {
	counterHandler := counter.NewCounterHandler(st)
	gaugeHandler := gauge.NewGaugeHandler(st)
	listHanlder := explorer.NewExplorerHandler(st)
	restHandler := rest.NewRestHandler(st)
	pingHandler := ping.NewPingHandler(st)
	router.Get("/", listHanlder.List)
	router.Get("/ping", pingHandler.Ping)
	router.Get("/features", pingHandler.Features)
	router.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		t := chi.URLParam(r, handlers.ParamTypeKey)
		switch t {
		case handlers.GaugeActionKey:
			gaugeHandler.Post(w, r)
		case handlers.CounterActionKey:
			counterHandler.Post(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	})
	router.Post("/update/{type}/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(GetErrorStatusCode(r))
	})
	router.Route("/value", func(r chi.Router) {
		r.Get("/gauge/{name}", gaugeHandler.Get)
		r.Get("/counter/{name}", counterHandler.Get)
	})
	router.Post("/update/", restHandler.Update)
	router.Post("/value/", restHandler.Value)
	router.Post("/updates/", restHandler.Updates)
}
