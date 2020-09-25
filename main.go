package main

import (
	"dota-predictor/app/config"
	"dota-predictor/app/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Set le port
func balanceTonPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func main() {
	//initialize the database
	config.InitDB()
	defer config.DB.Close()

	addr := balanceTonPort()
	r := mux.NewRouter().StrictSlash(true)

	log.Println("Starting server")
	handlers.HandleRequest(r)

	log.Fatal(http.ListenAndServe(addr, r))
}
