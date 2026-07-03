package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
	"url_shortner/internal/jwt"
	"url_shortner/internal/models"
	"url_shortner/internal/utils"
)

func (h *DB) LoginHandler(w http.ResponseWriter, req *http.Request) {
	user_credentials := models.Credentials{}
	err := json.NewDecoder(req.Body).Decode(&user_credentials)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}
	var hashedPass string
	var user_id int
	conn := h.Pool
	err = conn.QueryRow(context.Background(), "SELECT password,user_id FROM users WHERE email=$1", user_credentials.Email).Scan(&hashedPass, &user_id)
	if errors.Is(err, pgx.ErrNoRows) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Email/Password is not correct")
		return
	}
	err = utils.VerifyPassword(hashedPass, user_credentials.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("Invalid Password")
		return
	}
	jwt_token, err := jwt.Sign(user_id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Internal server error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(jwt_token)

}
