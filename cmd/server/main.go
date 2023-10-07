package main

import (
	"flag"

	"github.com/ekubyshin/metrics_agent/internal/server"
)

func main() {
	endpoint := flag.String("a", "localhost:8080", "endpoint address")
	flag.Parse()
	srv := server.NewServer(*endpoint)
	err := srv.Run()
	if err != nil {
		panic(err)
	}
}
