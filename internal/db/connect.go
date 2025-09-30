package db

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectToDatabase(ctx context.Context) (*pgx.Conn, error) {
	database := os.Getenv("DB_DATABASE")
	password := os.Getenv("DB_PASSWORD")
	username := os.Getenv("DB_USERNAME")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	schema := os.Getenv("DB_SCHEMA")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, url.QueryEscape(password), host, port, database, schema)
	db_conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return db_conn, nil
}
