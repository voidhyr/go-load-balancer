# go-load-balancer

An L7 HTTP load balancer built in Go with round-robin scheduling and active health checks.

## Features

- **Round-robin scheduling** — distributes requests evenly across backends
- **Active health checks** — pings backends every 5 seconds, automatically removes dead servers
- **Automatic recovery** — detects when a backend comes back online and restores it to the pool
- **Thread-safe** — per-server mutex ensures safe concurrent access
- **Timeout handling** — health checks have a 2 second timeout to avoid hanging

## Architecture

```
Client → Load Balancer (:8080) → Backend 1 (:8081)
                               → Backend 2 (:8082)
```

If a backend goes down, all traffic is rerouted to healthy backends automatically.
When it recovers, it rejoins the pool without restarting the balancer.

## Run it

**Start backend servers (separate terminals):**
```bash
go run server.go 8081
go run server.go 8082
```

**Start the load balancer:**
```bash
go run balancer.go
```

**Test round-robin:**
```bash
curl http://localhost:8080
curl http://localhost:8080
curl http://localhost:8080
curl http://localhost:8080
```

Expected output:
```
Hello from server on port 8081
Hello from server on port 8082
Hello from server on port 8081
Hello from server on port 8082
```

**Test health checks:**

Stop one backend with `Ctrl+C` — the balancer detects it within 5 seconds:
```
❌ Server DOWN: http://localhost:8081
```

All traffic automatically routes to the remaining healthy backend.
Restart the backend and it rejoins the pool:
```
✅ Server UP: http://localhost:8081
```

## Built with

- Go standard library
- `net/http/httputil` — reverse proxy
- `sync.Mutex` — thread-safe server state
- `net/http` — health check client with timeout

## Part of

Case Study: L7 HTTP Load Balancer
