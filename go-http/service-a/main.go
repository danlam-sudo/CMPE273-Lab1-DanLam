/*
CMPE 273 - Week 1 Lab 1: Your First Distributed System (Starter)
Dan Lam 011383814

Service A - Echo API
Hosted at http://127.0.0.1:8080

This is the provider service that provides the echo API.
This service is responsible for:
- /health: returns a JSON object with a status field set to "ok"
- /echo: returns a JSON object with the echo of the message
*/

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const serviceName = "echo-api"

// statusRecorder wraps http.ResponseWriter to capture status code.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(rec, r)
		log.Printf("service=%s method=%s path=%s status=%d latency_ms=%d",
			serviceName, r.Method, r.URL.Path, rec.status, time.Since(start).Milliseconds())
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	msg := r.URL.Query().Get("msg") // default "" if missing
	_ = json.NewEncoder(w).Encode(map[string]string{"echo": msg})
}

func main() {
	http.HandleFunc("/health", withLogging(health))
	http.HandleFunc("/echo", withLogging(echo))
	log.Printf("service=%s listening on 127.0.0.1:8080", serviceName)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
