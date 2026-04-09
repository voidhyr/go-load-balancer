package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

// Server represents a backend server
type Server struct {
	address string
	proxy   *httputil.ReverseProxy
}

// LoadBalancer holds all backend servers
type LoadBalancer struct {
	servers []*Server
	current int
	mu      sync.Mutex
}

// nextServer picks the next server (round robin)
func (lb *LoadBalancer) nextServer() *Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	server := lb.servers[lb.current]
	lb.current = (lb.current + 1) % len(lb.servers)
	return server
}

// ServeHTTP handles incoming requests
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.nextServer()
	fmt.Printf("Forwarding request to → %s\n", server.address)
	server.proxy.ServeHTTP(w, r)
}

func main() {
	// define your two backend servers
	addresses := []string{
		"http://localhost:8081",
		"http://localhost:8082",
	}

	// create server objects
	servers := []*Server{}
	for _, addr := range addresses {
		url, _ := url.Parse(addr)
		proxy := httputil.NewSingleHostReverseProxy(url)
		servers = append(servers, &Server{
			address: addr,
			proxy:   proxy,
		})
	}

	// create load balancer
	lb := &LoadBalancer{
		servers: servers,
	}

	fmt.Println("Load balancer running on port 8080")
	http.ListenAndServe(":8080", lb)
}
