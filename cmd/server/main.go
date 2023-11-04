package main

import (
	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/ekubyshin/metrics_agent/internal/logger"
	"github.com/ekubyshin/metrics_agent/internal/server"
)

func main() {
	cfg := config.AutoLoadServer()
	l, err := logger.NewLoggerFromEnv(&cfg)
	if err != nil {
		panic(err)
	}
	defer l.Sync() //nolint
	srv := server.NewServer(cfg, l)
	err = srv.Run()
	if err != nil {
		panic(err)
	}
}
