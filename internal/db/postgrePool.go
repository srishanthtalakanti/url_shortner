package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

func Config() (*pgxpool.Pool, error) {
	host := os.Getenv("host")
	port := os.Getenv("port")
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")

	psqlInfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user, password, host, port, dbname,
	)
	config, err := pgxpool.ParseConfig(psqlInfo)
	config.MinConns = 5
	config.MaxConns = 25

	config.MaxConnLifetime = 1 * time.Hour
	if err != nil {
		log.Println("Cant create a pool", err)
		return nil, err
	}
	return pgxpool.NewWithConfig(context.Background(), config)
}
