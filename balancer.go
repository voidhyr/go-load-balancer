package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// Server represents a backend server
type Server struct {
	address string
	proxy   *httputil.ReverseProxy
	alive   bool
	mu      sync.RWMutex
}

// isAlive safely checks if the server is alive
func (s *Server) isAlive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.alive
}

// setAlive safely sets the server's alive status
func (s *Server) setAlive(alive bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alive = alive
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

	for i := 0; i < len(lb.servers); i++ {
		server := lb.servers[lb.current]
		lb.current = (lb.current + 1) % len(lb.servers)

		if server.isAlive() {
			return server
		}
	}

	return nil // all servers are dead
}

// healthCheck pings a server to see if it is alive
func healthCheck(s *Server) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(s.address + "/health")
	if err != nil || resp.StatusCode != http.StatusOK {
		if s.isAlive() {
			fmt.Printf("❌ Server DOWN: %s\n", s.address)
		}
		s.setAlive(false)
		return
	}
	defer resp.Body.Close()

	if !s.isAlive() {
		fmt.Printf("✅ Server UP: %s\n", s.address)
	}
	s.setAlive(true)
}

// startHealthChecks runs health checks every 5 seconds
func startHealthChecks(lb *LoadBalancer) {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			start := time.Now()
			fmt.Printf("⏱ Cycle started at: %s\n",
				start.Format("15:04:05"))

			for _, server := range lb.servers {
				healthCheck(server)
			}

			fmt.Printf("⏱ Cycle ended at:   %s\n",
				time.Now().Format("15:04:05"))
		}
	}()
}

// ServeHTTP handles incoming requests
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.nextServer()

	// no alive servers
	if server == nil {
		http.Error(w, "No servers available", http.StatusServiceUnavailable)
		fmt.Println("⚠️  All servers are down!")
		return
	}

	fmt.Printf("Forwarding request → %s\n", server.address)
	server.proxy.ServeHTTP(w, r)
}

func main() {
	// define your three backend servers
	addresses := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	// create server objects
	servers := []*Server{}
	for _, addr := range addresses {
		url, err := url.Parse(addr)
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(url)
		servers = append(servers, &Server{
			address: addr,
			proxy:   proxy,
			alive:   true, // assume alive at start
		})
	}

	// create load balancer
	lb := &LoadBalancer{
		servers: servers,
	}

	startHealthChecks(lb)

	fmt.Println("Load balancer running on port 8080")
	http.ListenAndServe(":8080", lb)
}
