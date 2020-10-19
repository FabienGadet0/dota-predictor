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
		PredictProba        string
		PredictName         string
		ModelName           string
		RadiantTeam         string
		DireTeam            string
	}

	var data Data
	result := config.DB.Raw(`select p.*, g.radiant_team , g.dire_team, g.radiant_score, g.dire_score, g.start_date from prediction as p left join games as g on g.match_id = p.match_id order by inserted_date desc limit 25 offset ?`, offset).Scan(&data)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem retrieving the predictions from the database: " + result.Error.Error()})
		return
	}

	log.Println("/games-predicted for : " + r.Header.Get("access_token") + ".")
	json.NewEncoder(w).Encode(models.Response{Code: 0, Data: data})
}