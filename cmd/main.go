package main

import (
	"log"
	"net/http"
	"time"
	"url_shortner/internal/config"
	"url_shortner/internal/db"
	"url_shortner/internal/handlers"
	"url_shortner/internal/middleware"
	"url_shortner/internal/redis"
)

func main() {
	pool, err := db.Config()
	client := redis.NewRedisClient()
	if err != nil {
		log.Fatalf("Cant create pool")
	}
	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}
	handler := &handlers.DB{
		Pool:   pool,
		Client: client,
	}
	rl := &middleware.RateLimit{
		RedisClient: client,
		Limit:       10,
		Duration:    60 * time.Second,
	}
	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/", middleware.AuthMiddleware(rl.RateLimiter(handler.UrlHandler)))

	http.ListenAndServe(":8080", nil)
}
