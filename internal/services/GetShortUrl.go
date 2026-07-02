package services

import (
	"context"
	"github.com/jackc/pgx/v5"
	"url_shortner/internal/utils"
)

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
