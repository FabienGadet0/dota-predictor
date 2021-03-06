package handlers

import (
	"dota-predictor/app/config"
	"dota-predictor/app/helpers"
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

	if !isValidToken(w, r.Header.Get("access_token"), true, false) {
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

	if !isValidToken(w, r.Header.Get("access_token"), false, true) {
		return
	}

	mxline, err := strconv.Atoi(mux.Vars(r)["max-line"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem converting max-line" + err.Error()})
		return
	}
	if mxline < 5 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "Please enter a minimum lines > 5"})
		return
	}

	type Data struct {
		MatchID             int
		StartTime           *time.Time
		WinnerName          string
		PredictionIsCorrect bool
	}

	var data []Data
	result := config.DB.Raw(`select g.match_id, g.start_date , g.winner_name, p.predict_name = g.winner as prediction_is_correct 
	from games g 
	inner join prediction p on p.match_id = g.match_id and p.predict_name is not null 
	where p.model_name = 'main'
	order by g.start_date desc limit ?`, mxline).Scan(&data)

	if result.Error != nil || len(data) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem retrieving predictions from the database"})
		return
	}

	var score map[int]float64
	var countTrue float64

	score = make(map[int]float64)

		
	if len(data) < mxline {
		mxline = len(data)
	}

	for x, d := range data {
		if d.PredictionIsCorrect {
			countTrue++
		}
		switch v := x + 1; v {
		case 5:
			score[5] = countTrue / 5
		case mxline / 3:
			score[int(mxline / 3)] = countTrue / float64(mxline / 3)
		case mxline / 2:
			score[int(mxline / 2)] = countTrue / float64(mxline / 2)
		case mxline:
			score[int(mxline)] = countTrue / float64(mxline)
		}
	}

	log.Println("/model/score/{max-line} for : " + r.Header.Get("access_token") + ".")
	json.NewEncoder(w).Encode(models.Response{Code: 0, Data: score})
}

// @Summary Get last inserted date of a prediction
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Failure 423 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /model/last-run [get]
func getPredictionFromLastDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !isValidToken(w, r.Header.Get("access_token"), false, false) {
		return
	}

	type Data struct {
		InsertedDate *time.Time
	}

	var data Data
	result := config.DB.Raw(`select inserted_date from prediction order by inserted_date desc limit 1`).Scan(&data)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem retrieving the prediction from the database: " + result.Error.Error()})
		return
	}

	log.Println("/model/last-run for : " + r.Header.Get("access_token") + ".")
	json.NewEncoder(w).Encode(models.Response{Code: 0, Data: data})
}

// @Summary Get rows of predictions
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Failure 423 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /games-predicted [get]
func getPredictionsLive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !isValidToken(w, r.Header.Get("access_token"), false, true) {
		return
	}

	type Data struct {
		MatchID             int
		StartDate           *time.Time
		InsertedDate        *time.Time
		PredictProba        float64
		PredictName         string
		PredictTeam         string
		ModelName           string
		RadiantTeam         string
		DireTeam            string
	}

	var data []Data
	result := config.DB.Raw(`select p.*, 
	CASE WHEN p.predict_name = 'dire_team' THEN g.dire_team
		 WHEN p.predict_name = 'radiant_team' THEN g.radiant_team end
		  as predict_team,
				 g.radiant_team , g.dire_team, g.radiant_score, g.dire_score, g.start_date 
	from prediction as p 
	inner join games as g on g.match_id = p.match_id 
	where g.start_date >= now() - INTERVAL '3 hours'
	order by inserted_date desc 
	limit 100
	 `).Scan(&data)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem retrieving the predictions from the database: " + result.Error.Error()})
		return
	}

	log.Println("/games-predicted for : " + r.Header.Get("access_token") + ".")
	json.NewEncoder(w).Encode(models.Response{Code: 0, Data: data})
}

// @Summary Get rows of predictions
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Failure 423 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /games-predicted [get]
func getPredictions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query, ok := r.URL.Query()["page"]
	if !ok || len(query[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "Url parameter 'page' is missing."})
		return
	}

	offset, err := helpers.Pagination(query[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem generating the Pagination : " + err.Error()})
		return
	}

	if !isValidToken(w, r.Header.Get("access_token"), false, true) {
		return
	}

	type Data struct {
		MatchID             int
		StartDate           *time.Time
		InsertedDate        *time.Time
		PredictProba        float64
		PredictName         string
		PredictTeam         string
		ModelName           string
		RadiantTeam         string
		DireTeam            string
	}

	var data []Data
	result := config.DB.Raw(`select p.*, 
	CASE WHEN p.predict_name = 'dire_team' THEN g.dire_team
		 WHEN p.predict_name = 'radiant_team' THEN g.radiant_team end
		  as predict_team,
				 g.radiant_team , g.dire_team, g.radiant_score, g.dire_score, g.start_date 
	from prediction as p 
	left join games as g on g.match_id = p.match_id 
	order by inserted_date desc 
	limit 100 offset ?`, offset).Scan(&data)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem retrieving the predictions from the database: " + result.Error.Error()})
		return
	}

	log.Println("/games-predicted for : " + r.Header.Get("access_token") + ".")
	json.NewEncoder(w).Encode(models.Response{Code: 0, Data: data})
}

// @Summary Get all models name
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /model-name [get]
func getModelsNames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !isValidToken(w, r.Header.Get("access_token"), false, true) {
		return
	}

	type Data struct {
		ModelName string
	}
	var data []Data
	result := config.DB.Raw(`select distinct model_name from prediction`).Scan(&data)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem retrieving the models names from the database: " + result.Error.Error()})
		return
	}
	var d []string
	for _, x := range data {
		d = append(d, x.ModelName)
	}

	log.Println("/model-name for : " + r.Header.Get("access_token") + ".")
	json.NewEncoder(w).Encode(models.Response{Code: 0, Data: d})
}
