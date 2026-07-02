package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
	"url_shortner/internal/db"
	"url_shortner/internal/jwt"
	"url_shortner/internal/models"
	"url_shortner/internal/utils"
)

func LoginHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := db.ConnectDb()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Error connecting to database")
		return
	}
	defer conn.Close(context.Background())
	user_credentials := models.Credentials{}
	err = json.NewDecoder(req.Body).Decode(&user_credentials)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	}
	var hashedPass string
	var user_id int
	err = conn.QueryRow(context.Background(), "SELECT password,user_id FROM users WHERE email=$1", user_credentials.Email).Scan(&hashedPass, &user_id)
	if errors.Is(err, pgx.ErrNoRows) {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Email/Password is not correct")
		return
	}
	err = utils.VerifyPassword(hashedPass, user_credentials.Password)
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("Invalid Password")
		return
	}
	jwt_token, err := jwt.Sign(user_id)
	if err != nil {
		//what status code
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(jwt_token)

}
