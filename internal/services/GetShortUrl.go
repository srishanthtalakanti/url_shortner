package services

import (
	"context"
	"fmt"
	"log"
	"time"
	"url_shortner/internal/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func GetShortUrl(url string, conn *pgxpool.Pool, client *redis.Client, user_id int) (string, error) {
	//check if already exists if exsits return the short url
	var short_url string
	var err error
	err = conn.QueryRow(context.Background(), "SELECT short_code FROM urls WHERE long_url=$1", url).Scan(&short_url)
	switch err {
	case nil:
		return short_url, nil

	case pgx.ErrNoRows:

	default:
		return "", fmt.Errorf("lookup short url: %w", err)
	}
	//get next id in db

	var id int
	err = conn.QueryRow(context.Background(), "INSERT INTO urls (long_url,user_id) VALUES ($1,$2) RETURNING id", url, user_id).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert url: %w", err)
	}
	//hash function-> base 61 hashing to avoid /,+,=
	short_url = utils.HashFunction(id)
	//insert into redis too create a go routine and use synchronous channel for waiting
	done := make(chan bool, 1)
	go func() {
		if err := client.Set(context.Background(), short_url, url, time.Hour).Err(); err != nil {
			log.Println("cache set failed:", err)
		}
		done <- true
	}()
	//store the res in db
	_, err = conn.Exec(context.Background(), "UPDATE urls SET short_code=$1 WHERE id=$2", short_url, id)
	//return the res if cant store then return error
	if err != nil {
		return "", fmt.Errorf("update short code: %w", err)
	}
	<-done
	return short_url, err
	//return the res,error
}
