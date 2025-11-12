package main

import (
	"fmt"
	"log"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/healthz" {
		http.NotFound(w, req)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	n, err := w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		fmt.Println("Error writing response:", err)
		return
	}
	log.Printf("Wrote %d bytes to client\n", n)
}
