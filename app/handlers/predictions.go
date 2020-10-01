package handlers

import (
	"dota-predictor/app/config"
	"dota-predictor/app/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// @Summary Get prediction for specific match
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Failure 423 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /model/predict/{match-id} [get]
func getPrediction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !isValidToken(w, r.Header.Get("access_token"), true) {
		return
	}

	result := config.DB.Where("match_id = ?", mux.Vars(r)["match-id"]).First(&models.Prediction{})
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem retrieving the prediction from the database: " + result.Error.Error()})
		return
	}

	log.Println("/model/predict/{match-id} for : " + r.Header.Get("access_token") + ".")
	json.NewEncoder(w).Encode(models.Response{Code: 0, Data: result.Value})
}
