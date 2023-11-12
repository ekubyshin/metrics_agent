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
	gaugePostHandler := gauge.NewGaugePostHandler(st)
	counterPostHandler := counter.NewCounterPostHandler(st)
	gaugeGetHandler := gauge.NewGaugeGetHandler(st)
	counterGetHanlder := counter.NewCounterGetHandler(st)
	listHanlder := explorer.NewExplorerHandler(st)
	restHandler := rest.NewRestHandler(st)
	pingHandler := ping.NewPingHandler(st)
	router.Get(listHanlder.BaseURL(), listHanlder.ServeHTTP)
	router.Get(pingHandler.BaseURL(), pingHandler.ServeHTTP)
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
