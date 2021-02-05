//go:generate swag init

package main

import (
	"dota-predictor/app/config"
	"dota-predictor/app/handlers"
	"log"
	"net/http"
	"os"

	_ "dota-predictor/docs"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Set le port
func balanceTonPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4747"
		log.Println("INFO: No PORT environment variable detected, defaulting to " + port + ".")
	}
	return ":" + port
}

// @title Dota-Predictor API
// @version 1.0
// @description This is the documentation for the golang api of dota-predictor

// @contact.name API Support
// @contact.email pierre.saintsorny@gmail.com

// @BasePath /1.0
func main() {
	//initialize the database
	config.InitDB()
	defer config.DB.Close()

	addr := balanceTonPort()
	r := mux.NewRouter().StrictSlash(true)

	log.Println("Starting server.")

	handlers.HandleRequest(r)
    r.Use(mux.CORSMethodMiddleware(r))
	handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(addr, handler))
}
