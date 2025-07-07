package main

import (
	"log"

	"github.com/damiaoterto/vendas-proxy/internal/core"
)

func main() {
	proxy := core.NewProxy("0.0.0.0", 8080)
	if err := proxy.Listen(); err != nil {
		log.Fatal(err)
	}
}
