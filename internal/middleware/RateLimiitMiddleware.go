package middleware

import (
	"context"
	"fmt"
	"net/http"
)

func (rl *RateLimit) RateLimiter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userID, ok := req.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		key := fmt.Sprintf("rate_limit:%d", userID)

		count, err := rl.RedisClient.Incr(context.Background(), key).Result()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if count == 1 {
			if err := rl.RedisClient.Expire(context.Background(), key, rl.Duration).Err(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		if count > int64(rl.Limit) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, req)
	}
}
