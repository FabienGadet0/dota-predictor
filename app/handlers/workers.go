package handlers

import (
	"dota-predictor/app/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// @Summary Call a worker
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /predict/live [get]
func getLiveGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !isValidToken(w, r.Header.Get("access_token"), false, true) {
		return
	}

	response, err := http.Get(os.Getenv("WORKER_URL") + "/generate-live")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem calling the worker : " + err.Error()})
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem converting the data from the response : " + err.Error()})
		return
	}

	log.Println("/predict/live for : " + r.Header.Get("access_token") + ".")
	json.NewEncoder(w).Encode(models.Response{Code: 0, Data: body})
}

// @Summary Call a worker
// @Produce json
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /predict/recent-games [get]
func getRecentGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query, ok := r.URL.Query()["nb-days"]
	if !ok || len(query[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "Url parameter 'nb-days' is missing."})
		return
	}

	if !isValidToken(w, r.Header.Get("access_token"), false, true) {
		return
	}

	response, err := http.Get(os.Getenv("WORKER_URL") + "/generate-recent-games/" + query[0])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem calling the worker : " + err.Error()})
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem converting the data from the response : " + err.Error()})
		return
	}

	log.Println("/predict/recent-games for : " + r.Header.Get("access_token") + ".")
	json.NewEncoder(w).Encode(models.Response{Code: 0, Data: body})
}

// @Summary Call a worker
// @Produce json
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /predict/all [get]
func getAllGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !isValidToken(w, r.Header.Get("access_token"), false, true) {
		return
	}

	response, err := http.Get(os.Getenv("WORKER_URL") + "/predict-all")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem calling the worker : " + err.Error()})
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem converting the data from the response : " + err.Error()})
		return
	}

	log.Println("/predict/all for : " + r.Header.Get("access_token") + ".")
	json.NewEncoder(w).Encode(models.Response{Code: 0, Data: body})
}

// @Summary Call a worker
// @Produce json
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /train [get]
func getModelTrained(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !isValidToken(w, r.Header.Get("access_token"), false, true) {
		return
	}

	response, err := http.Get(os.Getenv("WORKER_URL") + "/train-all")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem calling the worker : " + err.Error()})
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem converting the data from the response : " + err.Error()})
		return
	}

	log.Println("/train for : " + r.Header.Get("access_token") + ".")
	json.NewEncoder(w).Encode(models.Response{Code: 0, Data: body})
}
