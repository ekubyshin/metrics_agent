package main

import (
	"github.com/ekubyshin/metrics_agent/internal/agent"
	"github.com/ekubyshin/metrics_agent/internal/config"
)

func main() {
	cfg := config.NewAgentConfig()
	metricsAgent := agent.NewMetricsAgent(cfg)
	metricsAgent.Start()
}
