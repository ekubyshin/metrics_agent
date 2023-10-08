package main

import (
	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/ekubyshin/metrics_agent/internal/server"
)

func main() {
	cfg := config.AutoLoad()
	srv := server.NewServer(cfg.Address.ToString())
	err := srv.Run()
	if err != nil {
		panic(err)
	}
}
