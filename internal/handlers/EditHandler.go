package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"url_shortner/internal/models"
)

func (h *DB) EditHandler(w http.ResponseWriter, req *http.Request) {
	url_body := models.EditUrl{}
	err := json.NewDecoder(req.Body).Decode(&url_body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}
	user_id := req.Context().Value("user_id")
	var id int
	conn := h.Pool
	conn.QueryRow(context.Background(), "SELECT user_id FROM urls WHERE short_code=$1", url_body.Short_code).Scan(&id)
	if id != user_id {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Not allowed to modify the url")
		return
	}
	_, err = conn.Exec(context.Background(), "UPDATE urls SET long_url=$1 WHERE short_code=$2", url_body.Long_url, url_body.Short_code)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Cant update the db")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode("Successfully updated the url")
}
