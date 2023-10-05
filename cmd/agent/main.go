package main

import (
	"time"

	"github.com/ekubyshin/metrics_agent/internal/agent"
)

func main() {
	metricsAgent := agent.NewMetricsAgent(10*time.Second, 2*time.Second)
	metricsAgent.Start()
}
