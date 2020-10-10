package handlers

import (
	"dota-predictor/app/config"
	"dota-predictor/app/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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

// @Summary Get prediction percentage on last x games
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /model/score/{max-line} [get]
func getPredictionPercentage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !isValidToken(w, r.Header.Get("access_token"), false) {
		return
	}

	mxline, err := strconv.Atoi(mux.Vars(r)["max-line"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem converting max-line" + err.Error()})
		return
	}
	if mxline != 5 && mxline != 10 && mxline != 20 && mxline != 50 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "Please enter 5, 10, 20 or 50 as a maximum number of lines"})
		return
	}

	type Data struct {
		MatchID             int
		StartTime           *time.Time
		WinnerName          string
		PredictionIsCorrect bool
	}

	var data []Data
	result := config.DB.Raw(`select g.match_id, g.start_date , g.winner_name, g.winner_name = p.predict_name as prediction_is_correct 
	from games g 
	inner join prediction p on p.match_id = g.match_id and p.predict_name is not null 
	order by g.start_date desc limit 50`).Scan(&data)

	if result.Error != nil || len(data) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem retrieving predictions from the database"})
		return
	}

	var score map[int]float64
	var countTrue float64

	score = make(map[int]float64)

	for x, d := range data {
		if d.PredictionIsCorrect {
			countTrue++
		}
		switch v := x + 1; v {
		case 5:
			score[5] = countTrue / 5
		case 10:
			score[10] = countTrue / 10
		case 20:
			score[20] = countTrue / 20
		case 50:
			score[50] = countTrue / 50
		}
	}
	if len(data) < mxline {
		score[len(data)] = countTrue / float64(len(data))
	}

	log.Println("/model/score/{max-line} for : " + r.Header.Get("access_token") + ".")
	json.NewEncoder(w).Encode(models.Response{Code: 0, Data: score})
}