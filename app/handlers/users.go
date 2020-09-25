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

func create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, err := helpers.TokenGenerator()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Code: -1, Message: "There was a problem generating the token"})
		return
	}

	call, err := strconv.Atoi(os.Getenv("MAX_NB_CALL_USER"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Code: -1, Message: "There was a problem converting MAX_NB_CALL_USER to integer"})
		return
	}

	var user = models.User{AccessToken: token, NBCallsLeft: call}

	result := config.DB.Create(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Code: -1, Message: "There was a problem inserting in database"})
		return
	}

	log.Println("New row for user : " + fmt.Sprintf("%v", result.Value))
	json.NewEncoder(w).Encode(models.Error{Code: 0, Data: result.Value})
}
