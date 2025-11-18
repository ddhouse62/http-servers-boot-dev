package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ddhouse62/http-servers-boot-dev/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Can't decode message")
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.dbQueries.LookupUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error looking up user")
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	check, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !check {
		log.Printf("Error could not verify password")
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if check {
		respondWithJSON(w, http.StatusOK, User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		})
	}
}
