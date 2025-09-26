package server

import (
	"rate-limiter/internal/db"
	ratelimit "rate-limiter/internal/limiter"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	DB      *pgx.Conn
	Queries *db.Queries
	Redis   *redis.Client
	Router  *chi.Mux
	Limiter *ratelimit.SlidingWindow
	Vars    *Vars
}

func New(
	db *pgx.Conn,
	queries *db.Queries,
	rdb *redis.Client,
	limiter *ratelimit.SlidingWindow,
) *Server {
	s := &Server{
		DB:      db,
		Queries: queries,
		Redis:   rdb,
		Router:  chi.NewRouter(),
		Limiter: limiter,
		Vars:    getVars(),
	}
	s.setRoutes()
	return s
}
