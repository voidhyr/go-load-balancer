package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Args[1] // pass port from terminal

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from server on port %s\n", port)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	fmt.Printf("Server running on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
