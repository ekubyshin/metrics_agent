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
	cfg := config.AutoLoadServer()
	l, err := logger.NewLoggerFromEnv(&cfg)
	if err != nil {
		panic(err)
	}
	defer l.Sync() //nolint
	var st storage.Storage[metrics.MetricsKey, metrics.Metrics]
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if cfg.DatabaseDSN != nil && *cfg.DatabaseDSN != "" {
		st, err = storage.NewDBStorage[metrics.MetricsKey, metrics.Metrics](ctx, &cfg)
		l.Info("use db")
		if err != nil {
			l.Info("db", err)
		}
	} else {
		memSt := storage.NewMemoryStorage[metrics.MetricsKey, metrics.Metrics]()
		if cfg.FileStoragePath != nil && *cfg.FileStoragePath != "" {
			st, err = storage.NewFileStorage(ctx, memSt, *cfg.FileStoragePath, *cfg.Restore, cfg.StoreDuration())
			if err != nil {
				panic(err)
			}
			l.Info("use filestorage")
		} else {
			l.Info("use memstorage")
			st = memSt
		}
	}
	srv := server.NewServer(cfg, l, st)
	err = srv.Run()
	if err != nil {
		panic(err)
	}
}
