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
		log.Fatalf("fail on load app config: %v", err)
	}

	mongo := database.NewMongoDB(config.MongoDB.URI)
	mongoConn, err := mongo.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer mongo.Disconnect()

	proxy := core.NewProxy(config, mongoConn)
	if err := proxy.Listen("0.0.0.0", 8080); err != nil {
		log.Fatal(err)
	}
}
