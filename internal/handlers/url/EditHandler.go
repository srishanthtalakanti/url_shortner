package url

import (
	"context"
	"encoding/json"
	"net/http"
	"url_shortner/internal/db"
	"url_shortner/internal/models"
)

func EditHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := db.ConnectDb()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Error connecting to database")
		return
	}
	defer conn.Close(context.Background())
	url_body := models.EditUrl{}
	err = json.NewDecoder(req.Body).Decode(&url_body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
	}
	user_id := req.Context().Value("user_id")
	var id int
	conn.QueryRow(context.Background(), "SELECT user_id FROM urls WHERE short_code=$1", url_body.Short_code).Scan(&id)
	if id != user_id {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Not allowed to modify the url")
		return
	}
	_, err = conn.Exec(context.Background(), "UPDATE urls SET long_url=$1 WHERE short_code=$2", url_body.Long_url, url_body.Short_code)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Cant update the db")
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode("Successfully updated the url")
}
