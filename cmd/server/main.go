package main

import (
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
)

func main() {
	router := http.NewServeMux()
	metricsHandler := handlers.NewMetricsHandler()
	router.Handle(metricsHandler.Route(), metricsHandler)
	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		panic(err)
	}
}
