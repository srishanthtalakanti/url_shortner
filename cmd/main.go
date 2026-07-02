package main

import (
	"net/http"
	"url_shortner/internal/handlers/auth"
	"url_shortner/internal/handlers/url"
	"url_shortner/internal/middleware"
)

func main() {
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/", middleware.AuthMiddleware(url.UrlHandler))

	http.ListenAndServe(":8080", nil)
}

//close the connections,dereferencing a nill pointer
