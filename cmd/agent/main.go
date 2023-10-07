package main

import (
	"flag"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/agent"
)

func main() {
	endpoint := flag.String("a", "localhost:8080", "endpoint address")
	reportInterval := flag.Int64("r", 10, "report interval")
	pollInterval := flag.Int64("p", 2, "poll interval")
	flag.Parse()
	metricsAgent := agent.NewMetricsAgent(time.Duration(*reportInterval)*time.Second, time.Duration(*pollInterval)*time.Second, *endpoint)
	metricsAgent.Start()
}
