package main

import (
	"context"

	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/ekubyshin/metrics_agent/internal/logger"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/server"
	"github.com/ekubyshin/metrics_agent/internal/storage"
)

func main() {
	cfg := config.NewServerConfig()
	l, err := logger.NewLoggerFromEnv(&cfg)
	if err != nil {
		panic(err)
	}
	defer l.Sync() //nolint
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	st, err := storage.AutoLoadStorage[metrics.MetricsKey, metrics.Metrics](ctx, cfg)
	if err != nil {
		panic(err)
	}
	srv := server.NewServer(cfg, l, st)
	err = srv.Run()
	if err != nil {
		panic(err)
	}
}
