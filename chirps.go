package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/ddhouse62/http-servers-boot-dev/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Can't decode message")
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140

	cleaned_body := profanityReplacement(params.Body)

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned_body,
		UserID: params.UserID,
	})
	if err != nil {
		log.Printf("Failed to create Chirp")
		respondWithError(w, http.StatusInternalServerError, "Couldn't Create Chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.CreatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

func profanityReplacement(msg string) string {
	words := strings.Split(msg, " ")

	profanity := []string{"kerfuffle", "sharbert", "fornax"}

	for i := 0; i < len(words); i++ {
		if slices.Contains(profanity, strings.ToLower(words[i])) {
			words[i] = "****"
		}
	}

	cleaned_words := strings.Join(words, " ")
	return cleaned_words
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error obtaining chirps")
		respondWithError(w, http.StatusInternalServerError, "Error obtaining requested chirps", err)
		return
	}
	response := []Chirp{}
	for _, chirp := range chirps {
		response = append(response, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Invalid Chirp ID: %s", r.PathValue("chirpID"))
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("Chirp not found")
		respondWithError(w, http.StatusNotFound, "Chirp Not Found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
