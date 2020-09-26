package handlers

import (
	"dota-predictor/app/config"
	"dota-predictor/app/helpers"
	"dota-predictor/app/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, err := helpers.TokenGenerator()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem generating the token: " + err.Error()})
		return
	}

	call, err := strconv.Atoi(os.Getenv("MAX_NB_CALL_USER"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem converting MAX_NB_CALL_USER to integer: " + err.Error()})
		return
	}

	var user = models.User{AccessToken: token, NBCallsLeft: call}

	result := config.DB.Create(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem inserting in database: " + result.Error.Error()})
		return
	}

	log.Println("New row for user : " + fmt.Sprintf("%v", result.Value))
	json.NewEncoder(w).Encode(models.Response{Code: 0, Data: result.Value})
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !isValidToken(w, r.Header.Get("access_token")) {
		return
	}

	result := config.DB.Where("access_token = ?", r.Header.Get("access_token")).First(&models.User{})
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem retrieving token from the database: " + result.Error.Error()})
		return
	}

	log.Println("New get for user : " + fmt.Sprintf("%v", result.Value))
	json.NewEncoder(w).Encode(models.Response{Code: 0, Data: result.Value})
}
