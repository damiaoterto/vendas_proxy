package main

import (
	"log"

	"github.com/damiaoterto/vendas-proxy/internal/config"
	"github.com/damiaoterto/vendas-proxy/internal/core"
	"github.com/damiaoterto/vendas-proxy/internal/database"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal("Fail on load app config")
	}

	_, err = database.ConnectFromURI(config.MongoDB.URI)
	if err != nil {
		log.Fatal(err)
	}

	proxy := core.NewProxy("0.0.0.0", 8080)
	if err := proxy.Listen(); err != nil {
		log.Fatal(err)
	}
}
