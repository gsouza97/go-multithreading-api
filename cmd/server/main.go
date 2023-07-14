package main

import (
	"net/http"

	"github.com/gsouza97/go-multithreading-api/internal/webserver/handlers"
)

func main() {
	mux := http.NewServeMux()
	cepHandler := handlers.NewCepHandler()
	mux.HandleFunc("/", cepHandler.GetCep)
	http.ListenAndServe(":3000", mux)
}
