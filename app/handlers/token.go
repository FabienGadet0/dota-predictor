package handlers

import (
	"dota-predictor/app/config"
	"dota-predictor/app/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

func isValidToken(w http.ResponseWriter, token string) bool {
	err := config.DB.Where("access_token = ?", token).Find(&models.User{}).Error
	if err == gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "Unknown token: " + err.Error()})
		return false
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem retrieving token from the database: " + err.Error()})
		return false
	}
	log.Println("User " + token + " just made a call to users.")
	return true
}
