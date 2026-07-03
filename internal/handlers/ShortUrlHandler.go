package handlers

import (
	"encoding/json"
	"net/http"
	"url_shortner/internal/models"
	"url_shortner/internal/services"
)

func (h *DB) ShortUrlHandler(w http.ResponseWriter, req *http.Request) {
	var body models.ShortenUrl
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Invalid/Empty url")
		return
	}
	url := body.URL
	conn := h.Pool
	res, err := services.GetShortUrl(url, conn, h.Client, req.Context().Value("user_id").(int))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
	}
}
