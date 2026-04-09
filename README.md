# go-load-balancer

An L7 HTTP load balancer built in Go with round-robin scheduling.

## How it works

Incoming requests on port 8080 are distributed across backend servers using round-robin. A mutex ensures safe concurrent access to the server counter.

## Run it

**Start backend servers:**
```bash
go run server.go 8081
go run server.go 8082
```

**Start the load balancer:**
```bash
go run balancer.go
```

**Test it:**
```bash
curl http://localhost:8080
curl http://localhost:8080
curl http://localhost:8080
curl http://localhost:8080
```

## Architecture

## Built with
- Go standard library
- `net/http/httputil` reverse proxy
