/*
CMPE 273 - Week 1 Lab 1: Your First Distributed System (Starter)
Dan Lam 011383814

Service B - Client
Hosted at http://127.0.0.1:8081

This is the client service that calls the provider service.
This service is responsible for:
- /health: returns a JSON object with a status field set to "ok"
- /call-echo: calls the provider service and returns a JSON object with the response from the provider service
*/

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

const serviceName = "client"
const providerBase = "http://127.0.0.1:8080"

// client is the HTTP client that is used to call the provider service.
// It is configured with a timeout of 2 seconds.
var client = &http.Client{Timeout: 2 * time.Second}

// statusRecorder wraps http.ResponseWriter to capture status code.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Latency logging middleware
// It logs the service name, method, path, status code, and latency.
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

func callEcho(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	msg := r.URL.Query().Get("msg")

	u, _ := url.Parse(providerBase + "/echo")
	q := u.Query()
	q.Set("msg", msg)
	u.RawQuery = q.Encode()

	resp, err := client.Get(u.String())
	if err != nil {
		latency := time.Since(start).Milliseconds()
		log.Printf("service=%s endpoint=/call-echo status=503 latency_ms=%d error=%q",
			serviceName, latency, err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":   "provider unavailable",
			"details": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var providerResp struct {
		Echo string `json:"echo"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&providerResp); err != nil {
		latency := time.Since(start).Milliseconds()
		log.Printf("service=%s endpoint=/call-echo status=503 latency_ms=%d error=%q",
			serviceName, latency, err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":   "provider unavailable",
			"details": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"consumer": map[string]string{"msg": msg},
		"provider": map[string]string{"echo": providerResp.Echo},
	})
}

func main() {
	http.HandleFunc("/health", withLogging(health))
	http.HandleFunc("/call-echo", withLogging(callEcho))
	log.Printf("service=%s listening on 127.0.0.1:8081", serviceName)
	log.Fatal(http.ListenAndServe("127.0.0.1:8081", nil))
}
