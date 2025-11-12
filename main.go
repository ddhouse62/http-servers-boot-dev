package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const filepathRoot = "."
	const port = "8080"
	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)

	mux.HandleFunc("/reset", apiCfg.handlerReset)

	mux.HandleFunc("/healthz", handlerReadiness)
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Load()

	statusText := fmt.Sprintf("Hits: %d", hits)

	if r.URL.Path != "/metrics" {
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
