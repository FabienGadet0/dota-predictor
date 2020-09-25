package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// HandleRequest gere toute les routes du serveur HTTP
func HandleRequest(router *mux.Router) {
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/user/create", create).Methods("POST")
}

//Base route
func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json;charset=UTF-8")
	fmt.Fprintf(w, "<h1>Hi there, welcome to the dota-predictor api !</h1>")
}
