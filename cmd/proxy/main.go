package main

import (
	"log"

	"github.com/damiaoterto/vendas-proxy/internal/config"
	"github.com/damiaoterto/vendas-proxy/internal/core"
)

func main() {
	_, err := config.Load()
	if err != nil {
		log.Fatalf("Fail on load app config")
	}

	proxy := core.NewProxy("0.0.0.0", 8080)
	if err := proxy.Listen(); err != nil {
		log.Fatal(err)
	}
}
