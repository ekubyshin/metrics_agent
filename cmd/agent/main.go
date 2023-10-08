package main

import (
	"flag"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/agent"
	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/ekubyshin/metrics_agent/internal/utils"
)

func main() {
	cfg := InitConfig()
	metricsAgent := agent.NewMetricsAgent(cfg)
	metricsAgent.Start()
}

func InitConfig() *config.Config {
	cfg := config.NewConfig()
	endpoint := flag.String("a", "localhost:8080", "endpoint address")
	reportInterval := flag.Int("r", 10, "report interval")
	pollInterval := flag.Int("p", 2, "poll interval")
	flag.Parse()
	if cfg.Address == nil {
		cfg.Address = &config.Address{}
		cfg.Address.UnmarshalText([]byte(*endpoint))
	}
	if cfg.PollInterval == nil {
		cfg.PollInterval = utils.ToPointer[time.Duration](utils.IntToDuration(*pollInterval))
	}
	if cfg.ReportInterval == nil {
		cfg.ReportInterval = utils.ToPointer[time.Duration](utils.IntToDuration(*reportInterval))
	}
	return cfg
}
