package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func handlerChirpValidator(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
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
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleaned_body,
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
