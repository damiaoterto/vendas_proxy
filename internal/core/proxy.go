package core

import (
	"fmt"
	"net/http"
	"strings"
)

type Proxy struct {
	mux  *http.ServeMux
	addr string
	port int
}

func NewProxy(addr string, port int) *Proxy {
	mux := http.NewServeMux()
	return &Proxy{addr: addr, port: port, mux: mux}
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
func (p *Proxy) Listen() error {
	p.mux.HandleFunc("/", p.proxyHandler)

	addr := fmt.Sprintf("%s:%d", p.addr, p.port)

	if err := http.ListenAndServe(addr, p.mux); err != nil {
		return fmt.Errorf("error on start http server: %v", err)
	}

	return nil
}
