package main

import (
	"log"
	"net/http"

)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)

	if r.URL.Path != "/api/reset" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	n, err := w.Write([]byte("File Server Hits reset to 0"))
	if err != nil {
		log.Println("Error writing response:", err)
		return
	}
	log.Printf("Wrote %d bytes to client\n", n)

}
