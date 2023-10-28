package main

import (
	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/ekubyshin/metrics_agent/internal/logger"
	"github.com/ekubyshin/metrics_agent/internal/server"
)

func main() {
	cfg := config.AutoLoad()
	var l logger.Logger
	var err error
	if cfg.Env == "production" {
		l, err = logger.NewProductionLogger()
	} else {
		l, err = logger.NewDevelopmentLogger()
	}
	if err != nil {
		panic(err)
	}
	defer l.Sync() //nolint
	srv := server.NewServer(cfg.Address, l)
	err = srv.Run()
	if err != nil {
		panic(err)
	}
}
