package url

import (
	"encoding/json"
	"net/http"
)

func UrlHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		RedirectHandler(w, req) //converts from short url to long url
	case http.MethodPost:
		ShortUrlHandler(w, req) //converts from long url to short url
	case http.MethodPatch:
		EditHandler(w, req)
	default:
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Invalid Http Method")
	}

}
