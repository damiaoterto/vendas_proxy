package core

import (
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Proxy struct {
	mux   *http.ServeMux
	mongo *mongo.Client
}

func NewProxy(mongo *mongo.Client) *Proxy {
	mux := http.NewServeMux()
	return &Proxy{mongo: mongo, mux: mux}
}

func (p Proxy) proxyHandler(w http.ResponseWriter, r *http.Request) {
	method := strings.ToLower(r.Method)
	if method != "get" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hostParts := strings.Split(r.Host, ".")
	if len(hostParts) < 2 {
		http.Error(w, "Invalid host", http.StatusBadRequest)
		return
	}
	subdomain := hostParts[0]

	fmt.Println(subdomain)
}
func (p *Proxy) Listen(addr string, port int) error {
	p.mux.HandleFunc("/", p.proxyHandler)

	addr = fmt.Sprintf("%s:%d", addr, port)

	if err := http.ListenAndServe(addr, p.mux); err != nil {
		return fmt.Errorf("error on start http server: %v", err)
	}

	return nil
}
