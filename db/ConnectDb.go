package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

func ConnectDb() (*pgx.Conn, error) {
	host := os.Getenv("host")
	port := os.Getenv("port")
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")

	psqlInfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user, password, host, port, dbname,
	)
	conn, err := pgx.Connect(context.Background(), psqlInfo)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
