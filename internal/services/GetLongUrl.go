package services

import (
	"context"
	"github.com/jackc/pgx/v5"
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
