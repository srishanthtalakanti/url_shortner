package services

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"time"
)

func GetLongUrl(short_code string, conn *pgxpool.Pool, client *redis.Client) (string, error) {
	var long_url string
	var err error
	long_url, err = client.Get(context.Background(), short_code).Result()
	if err != nil {
		err = conn.QueryRow(context.Background(), "SELECT long_url FROM urls WHERE short_code=$1", short_code).Scan(&long_url)
		if err != nil {
			return "", err
		}
		client.Set(context.Background(), short_code, long_url, time.Hour)
	}
	return long_url, nil

}
