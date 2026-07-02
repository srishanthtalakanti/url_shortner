package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"url_shortner/internal/jwt"
	"url_shortner/internal/models"
	"url_shortner/internal/utils"
)

func (h *DB) RegisterHandler(w http.ResponseWriter, req *http.Request) {
	//get the credentials

	user_credentials := models.Credentials{}
	err := json.NewDecoder(req.Body).Decode(&user_credentials)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Invalid/Empty credentials")
		return
	}
	hashedPass, err := utils.HashPassword(user_credentials.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Yoo")
		return
	}
	//insert into users table and get the user_id
	var user_id int
	err = h.Pool.QueryRow(context.Background(), "INSERT INTO users (email,password) VALUES ($1,$2) RETURNING user_id", user_credentials.Email, hashedPass).Scan(&user_id)
	if err != nil {
		//what status code?
		w.WriteHeader(500)
		log.Println(err)
		json.NewEncoder(w).Encode("	Internal Server Error")
		return
	}
	//get the jwt
	jwt_token, err := jwt.Sign(user_id)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Internal Server Error")
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(jwt_token)
	//if err show the error

}
