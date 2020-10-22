package handlers

import (
	"dota-predictor/app/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	routes []string
)

func gorillaWalkFn(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	path, _ := route.GetPathTemplate()
	if path == "/"+os.Getenv("VERSION") {
		return nil
	}
	routes = append(routes, path)
	return nil
}

// HandleRequest gere toute les routes du serveur HTTP
func HandleRequest(router *mux.Router) {
	// Subrouter to get the url match with the api version
	r := router.PathPrefix("/" + os.Getenv("VERSION")).Subrouter().StrictSlash(true)

	r.HandleFunc("", index).Methods("GET")
	r.HandleFunc("/users/personnenetrouverajamaismaroutedecreationdutilisateur", createUser).Methods("POST")
	r.HandleFunc("/users/stats", getUser).Methods("GET")
	r.HandleFunc("/list-routes", listRoutes).Methods("GET")
	r.HandleFunc("/model/predict/{match-id}", getPrediction).Methods("GET")
	r.HandleFunc("/model/score/{max-line}", getPredictionPercentage).Methods("GET")
	r.HandleFunc("/model/last-run", getPredictionFromLastDate).Methods("GET")
	r.HandleFunc("/games-predicted", getPredictions).Methods("GET")
	r.HandleFunc("/games-predicted-live", getPredictionsLive).Methods("GET")
	r.HandleFunc("/model-name", getModelsNames).Methods("GET")

	r.HandleFunc("/predict/live", getLiveGames).Methods("GET")
	r.HandleFunc("/predict/recent-games", getRecentGames).Methods("GET")
	r.HandleFunc("/predict/all", getAllGames).Methods("GET")
	r.HandleFunc("/train", getModelTrained).Methods("GET")

	// Documentation
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// Get all routes
	err := router.Walk(gorillaWalkFn)
	if err != nil {
		log.Fatal(err)
	}
}

// @Summary Index of the api
// @Produce json
// @Success 200 {object} models.Response
// @Router / [get]
func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json;charset=UTF-8")

	fmt.Fprintf(w, "<h1>Hi there, welcome to the dota-predictor api !</h1>")

	json.NewEncoder(w).Encode(models.Response{Code: 0})
}

// @Summary List all routes defined in the api
// @Produce json
// @Success 200 {object} string
// @Router /list-routes [get]
func listRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json;charset=UTF-8")

	if !isValidToken(w, r.Header.Get("access_token"), false, false) || (*r).Method == "OPTIONS" {
		return
	}

	json.NewEncoder(w).Encode(routes)
}
