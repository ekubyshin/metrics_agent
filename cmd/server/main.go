package main

import (
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
)

func main() {
	router := http.NewServeMux()
	metricsHandler := handlers.NewMetricsHandler()
	router.Handle(metricsHandler.Route(), metricsHandler)
	http.ListenAndServe("127.0.0.1:8080", router)
}
