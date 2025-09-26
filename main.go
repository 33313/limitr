package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"rate-limiter/internal/db"
	"rate-limiter/internal/limiter"
	"rate-limiter/internal/server"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	godotenv.Load()

	var conn *pgx.Conn
	ctx := context.Background()
	queries, err := db.ConnectToDatabase(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})

	l := limiter.New(rdb)

	s := server.New(conn, queries, rdb, l)

	http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), s.Router)
}
