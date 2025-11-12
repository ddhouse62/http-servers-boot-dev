package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Load()

	statusText := fmt.Sprintf("Hits: %d", hits)

	if r.URL.Path != "/api/metrics" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	n, err := w.Write([]byte(statusText))
	if err != nil {
		log.Println("Error writing response:", err)
		return
	}
	log.Printf("Wrote %d bytes to client\n", n)
}
