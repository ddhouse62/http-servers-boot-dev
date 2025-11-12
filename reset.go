package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")

	if platform != "dev" {
		log.Printf("Access forbidden - Not in Dev environment")
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}

	cfg.dbQueries.ResetUser(r.Context())

	cfg.fileserverHits.Store(0)

	if r.URL.Path != "/admin/reset" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	n, err := w.Write([]byte("Server reset"))
	if err != nil {
		log.Println("Error writing response:", err)
		return
	}
	log.Printf("Wrote %d bytes to client\n", n)

}
