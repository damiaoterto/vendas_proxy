package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/damiaoterto/vendas-proxy/internal/config"
	"github.com/damiaoterto/vendas-proxy/internal/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Proxy struct {
	mongo  *mongo.Client
	config *config.AppConfig
}

type RouteData struct {
	Subdomain string `json:"subdomain"`
	TargetURL string `json:"target_url"`
}

func NewProxy(config *config.AppConfig, mongo *mongo.Client) *Proxy {
	return &Proxy{config: config, mongo: mongo}
}

func (p Proxy) CreateNewRoute(w http.ResponseWriter, r *http.Request) {
	var routeData RouteData
	if err := json.NewDecoder(r.Body).Decode(&routeData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	collection := p.mongo.Database(p.config.MongoDB.Database).Collection(p.config.MongoDB.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	route := model.Route{
		Subdomain: routeData.Subdomain,
		TargetURL: routeData.TargetURL,
		IsActive:  true,
	}

	_, err := collection.InsertOne(ctx, route)
	if err != nil {
		log.Printf("Error inserting route into MongoDB: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message":    "Route created successfully",
		"subdomain":  routeData.Subdomain,
		"target_url": routeData.TargetURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
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
	r := mux.NewRouter()

	r.HandleFunc("/", p.proxyHandler).Methods("GET")
	r.HandleFunc("/routes", p.CreateNewRoute).Methods("POST")

	addr = fmt.Sprintf("%s:%d", addr, port)

	if err := http.ListenAndServe(addr, r); err != nil {
		return fmt.Errorf("error on start http server: %v", err)
	}

	return nil
}
