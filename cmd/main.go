package main

import (
	"log"
	"net/http"
	"url_shortner/internal/db"
	"url_shortner/internal/handlers"
	"url_shortner/internal/middleware"
)

func main() {
	pool, err := db.Config()
	if err != nil {
		log.Fatalf("Cant create pool")
	}
	h := handlers.DB{
		Pool: pool,
	}
	http.HandleFunc("/register", h.RegisterHandler)
	http.HandleFunc("/login", h.LoginHandler)
	http.HandleFunc("/", middleware.AuthMiddleware(h.UrlHandler))

	http.ListenAndServe(":8080", nil)
}

//close the connections,dereferencing a nill pointer
