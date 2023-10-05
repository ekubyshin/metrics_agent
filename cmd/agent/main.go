package main

import (
	"github.com/ekubyshin/metrics_agent/internal/agent"
)

func main() {
	metricsAgent := agent.NewMetricsAgent(2, 10)
	metricsAgent.Start()
}
