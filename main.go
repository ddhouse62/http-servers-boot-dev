package main

import (
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"database/sql"

	"github.com/ddhouse62/http-servers-boot-dev/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening sql connection: %s", err)
	}

	const filepathRoot = "."
	const port = "8080"
	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      database.New(db),
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirps)
	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	jwtSecret      string
}
