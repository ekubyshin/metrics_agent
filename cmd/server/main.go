package main

import (
	"flag"

	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/ekubyshin/metrics_agent/internal/server"
)

func main() {
	cfg := InitConfig()
	srv := server.NewServer(cfg.Address.ToString())
	err := srv.Run()
	if err != nil {
		panic(err)
	}
}

func InitConfig() *config.Config {
	cfg := config.NewConfig()
	endpoint := flag.String("a", "localhost:8080", "endpoint address")
	flag.Parse()
	if cfg.Address == nil {
		cfg.Address = &config.Address{}
		cfg.Address.UnmarshalText([]byte(*endpoint))
	}
	return cfg
}
