package handlers

import (
	"dota-predictor/app/config"
	"dota-predictor/app/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// check the validity of the token and decrement a call if bool is true
func isValidToken(w http.ResponseWriter, token string, decrementCall bool, isLevelGranted bool) bool {
	var user models.Users
	enableCors(&w)
	err := config.DB.Where("access_token = ?", token).Find(&user).Error
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

	if decrementCall {
		if user.NBCallsLeft == 0 {
			w.WriteHeader(http.StatusLocked)
			json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "No call left available."})
			return false
		}
		user.NBCallsLeft--
		config.DB.Save(&user)
	}

	if isLevelGranted {
		if user.GrantLvl != 1 {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "Access forbidden for this user (grant_lvl = 0)."})
			return false
		}
	}

	log.Println("User " + token + " verified.")
	return true
}
