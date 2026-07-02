package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"url_shortner/internal/jwt"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		cliams, err := jwt.ValidateToken(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusNetworkAuthenticationRequired)
			json.NewEncoder(w).Encode("Unauthenticated User")
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", cliams.UserId)
		next(w, r.WithContext(ctx))
	})
}
