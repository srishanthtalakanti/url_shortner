package url

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"url_shortner/internal/db"
	"url_shortner/internal/services"
)

func RedirectHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := db.ConnectDb()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Error connecting to database")
		return
	}
	defer conn.Close(context.Background())
	id := strings.TrimPrefix(req.URL.Path, "/")
	//get the url from the query
	res, err := services.GetLongUrl(id, conn) //has the short url now get the long url
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
