package main

import (
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	serve := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	serve.ListenAndServe()

}
