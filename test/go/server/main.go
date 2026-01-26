// main.go
// Simple HTTP server for Docker container
// Listens on port 8080 and responds with a plain text message
package main

import (
	"fmt"
	"net/http"
	"time"
)

const (
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 60 * time.Second
)

func handler(w http.ResponseWriter, _ *http.Request) { //nolint:unused // _ is unused as this is a simple handler
	w.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprintln(w, "Hello from Go HTTP server!")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", handler)
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
