package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"url_shortner/internal/services"
)

func (h *DB) RedirectHandler(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/")
	//get the url from the query
	res, err := services.GetLongUrl(id, h.Pool) //has the short url now get the long url
	//handle the error case
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)

	} else {
		w.WriteHeader(http.StatusPermanentRedirect)
		http.Redirect(w, req, res, http.StatusPermanentRedirect)
		//status code -> redirect 301
		//return the res(set the headers,status code,etc)
	}

}
