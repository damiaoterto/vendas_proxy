package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/damiaoterto/vendas-proxy/internal/config"
	"github.com/damiaoterto/vendas-proxy/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Proxy struct {
	mux    *http.ServeMux
	mongo  *mongo.Client
	config *config.AppConfig
}

func NewProxy(config *config.AppConfig, mongo *mongo.Client) *Proxy {
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

	log.Printf("Start proxy with subdomain %s", subdomain)

	collection := p.mongo.Database(p.config.MongoDB.Database).Collection(p.config.MongoDB.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var route model.Route
	filter := bson.M{"subdomain": subdomain}

	err := collection.FindOne(ctx, filter).Decode(&route)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Subdomain not found", http.StatusNotFound)
			return
		}

		log.Printf("Erro ao consultar o MongoDB: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	targetURL, err := url.Parse(route.TargetURL)
	if err != nil {
		log.Printf("Invalid target url to domain %s: %s", subdomain, route.TargetURL)
		http.Error(w, "Invalid route config", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = targetURL.Host
	}

	log.Printf("Redirecionando: %s (subdomínio: %s) -> %s", r.Host, subdomain, targetURL.String())
	proxy.ServeHTTP(w, r)
}

func (p *Proxy) Listen(addr string, port int) error {
	p.mux.HandleFunc("/", p.proxyHandler)

	addr = fmt.Sprintf("%s:%d", addr, port)

	if err := http.ListenAndServe(addr, p.mux); err != nil {
		return fmt.Errorf("error on start http server: %v", err)
	}

	return nil
}
