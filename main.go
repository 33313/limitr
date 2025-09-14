package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"rate-limiter/internal/api"
	"rate-limiter/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

var (
	ctx     = context.Background()
	conn    *pgx.Conn
	queries *db.Queries
)

type ApiKey struct {
	Id             string `json:"id"`
	Key            string `json:"hashed_key"`
	LimitPerMinute int    `json:"limit_per_minute"`
}

func connectToDatabase() error {
	database := os.Getenv("DB_DATABASE")
	password := os.Getenv("DB_PASSWORD")
	username := os.Getenv("DB_USERNAME")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	schema := os.Getenv("DB_SCHEMA")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, url.QueryEscape(password), host, port, database, schema)
	dbConn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return err
	}
	conn = dbConn
	queries = db.New(dbConn)
	return nil
}

func main() {
	godotenv.Load()
	if err := connectToDatabase(); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := conn.Ping(ctx); err != nil {
			http.Error(w, "Failed to ping database", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/keys", func(w http.ResponseWriter, r *http.Request) {
		limit_per_minute := 60
		key, err := api.GenerateAPIKey()
		if err != nil {
			http.Error(w, "Failed to generate key", http.StatusInternalServerError)
			return
		}

		row, err := queries.CreateKey(ctx, db.CreateKeyParams{
			HashedKey: api.HashAPIKey(key),
		})
		if err != nil {
			log.Print(err)
			http.Error(w, "Failed to create key", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		data := ApiKey{
			Id:             row.ID.String(),
			Key:            key,
			LimitPerMinute: limit_per_minute,
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	})

	r.Get("/keys", func(w http.ResponseWriter, r *http.Request) {
		rows, err := queries.ListKeys(ctx)
		if err != nil {
			http.Error(w, "Failed to list keys", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rows)
	})

	r.Delete("/keys/{key_id}", func(w http.ResponseWriter, r *http.Request) {
		key_uuid_string := chi.URLParam(r, "key_id")
		var key_uuid pgtype.UUID
		if err := key_uuid.Scan(key_uuid_string); err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}
		if err := queries.DeleteKey(ctx, key_uuid); err != nil {
			http.Error(w, "Failed to delete key", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	http.ListenAndServe(":3000", r)
}
