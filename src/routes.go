package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//Gere toute les routes du serveur HTTP
func handleRequest(router *mux.Router) {
	router.HandleFunc("/", index).Methods("GET")
}

//Base route
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hi there, welcome to the dota-predictor api !</h1>")
}
