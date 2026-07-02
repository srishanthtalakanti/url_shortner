package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *DB) UrlHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		h.RedirectHandler(w, req) //converts from short url to long url
	case http.MethodPost:
		h.ShortUrlHandler(w, req) //converts from long url to short url
	case http.MethodPatch:
		h.EditHandler(w, req)
	default:
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Invalid Http Method")
	}

}
