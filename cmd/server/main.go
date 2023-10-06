package main

import (
	"github.com/ekubyshin/metrics_agent/internal/server"
)

func main() {
	srv := server.NewServer()
	err := srv.Run()
	if err != nil {
		panic(err)
	}
}
