package handlers

import (
	"dota-predictor/app/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// HandleRequest gere toute les routes du serveur HTTP
func HandleRequest(router *mux.Router) {
	router.HandleFunc("/v1/", index).Methods("GET")
	router.HandleFunc("/v1/users/personnenetrouverajamaismaroutedecreationdutilisateur", createUser).Methods("POST")
	router.HandleFunc("/v1/users/stats", getUser).Methods("GET")
}

//Base route
func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json;charset=UTF-8")

	fmt.Fprintf(w, "<h1>Hi there, welcome to the dota-predictor api !</h1>")

	json.NewEncoder(w).Encode(models.Response{Code: 0})
}
