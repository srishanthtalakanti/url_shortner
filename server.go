package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strings"
	"url_shortner/db"
	"url_shortner/jwt"
	"url_shortner/middleware"
	"url_shortner/models"
	"url_shortner/utils"
)

func GetLongUrl(id string, conn *pgx.Conn) (string, error) {
	var long_url string
	var err error
	err = conn.QueryRow(context.Background(), "SELECT long_url FROM urls WHERE short_code=$1", id).Scan(&long_url)
	if err != nil {
		return "", err
	}
	return long_url, nil

}

func GetShortUrl(url string, conn *pgx.Conn, user_id int) (string, error) {
	//check if already exists if exsits return the short url
	var short_url string
	var err error
	err = conn.QueryRow(context.Background(), "SELECT short_code FROM urls WHERE long_url=$1", url).Scan(&short_url)
	if err == nil {
		return short_url, nil
	}
	//get next id in db
	var id int
	err = conn.QueryRow(context.Background(), "INSERT INTO urls (long_url,user_id) VALUES ($1,$2) RETURNING id,user_id", url, user_id).Scan(&id)
	if err != nil {
		return "", err
	}
	//hash function-> base 61 hashing to avoid /,+,=
	hashValue := utils.HashFunction(id)
	//store the res in db
	_, err = conn.Exec(context.Background(), "UPDATE urls SET short_code=$1 WHERE id=$2", hashValue, id)
	//return the res if cant store then return error
	if err != nil {
		return "", err
	}
	return hashValue, err
	//return the res,error
}

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
	res, err := GetLongUrl(id, conn) //has the short url now get the long url
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
func ShortUrlHandler(w http.ResponseWriter, req *http.Request) {
	var body models.ShortenUrl
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Invalid/Empty url")
		return
	}
	url := body.URL
	conn, err := db.ConnectDb()

	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Error connecting to database")
		return
	}
	defer conn.Close(context.Background())
	res, err := GetShortUrl(url, conn, req.Context().Value("user_id").(int))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
	}
}

func EditHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := db.ConnectDb()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Error connecting to database")
		return
	}
	defer conn.Close(context.Background())
	url_body := models.EditUrl{}
	err = json.NewDecoder(req.Body).Decode(&url_body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
	}
	user_id := req.Context().Value("user_id")
	var id int
	conn.QueryRow(context.Background(), "SELECT user_id FROM urls WHERE short_code=$1", url_body.Short_code).Scan(&id)
	if id != user_id {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Not allowed to modify the url")
		fmt.Println(id)
		fmt.Println(user_id)
		return
	}
	_, err = conn.Exec(context.Background(), "UPDATE urls SET long_url=$1 WHERE short_code=$2", url_body.Long_url, url_body.Short_code)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Cant update the db")
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode("Successfully updated the url")
}
func urlHandler(w http.ResponseWriter, req *http.Request) {
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

func registerHandler(w http.ResponseWriter, req *http.Request) {
	//get the credentials
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
		json.NewEncoder(w).Encode("Invalid/Empty credentials")
		return
	}
	hashedPass, err := utils.HashPassword(user_credentials.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal Server Error")
		return
	}
	//insert into users table and get the user_id
	var user_id int
	err = conn.QueryRow(context.Background(), "INSERT INTO users (email,password) VALUES ($1,$2) RETURNING user_id", user_credentials.Email, hashedPass).Scan(&user_id)
	if err != nil {
		//what status code?
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}
	//get the jwt
	jwt_token, err := jwt.Sign(user_id)
	if err != nil {
		//what status code
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(jwt_token)
	//if err show the error

}
func loginHandler(w http.ResponseWriter, req *http.Request) {
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
func main() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/", middleware.AuthMiddleware(urlHandler))

	http.ListenAndServe(":8080", nil)
}

//close the connections,dereferencing a nill pointer
