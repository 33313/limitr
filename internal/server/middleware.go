package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"rate-limiter/internal/keys"
	"strconv"
	"time"
)

const (
	MIDDLEWARE_PLAN_KEY = "limitr_plan:%v"
)

func (s *Server) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api_key := r.Header.Get("x-api-key")
		if api_key == "" {
			http.Error(w, "No API key provided", http.StatusUnauthorized)
			return
		}

		hashed_key := keys.HashKey(api_key)

		var (
			max_reqs       int
			window_seconds int
		)

		max_reqs, window_seconds, err := s.getPlan(w, r, hashed_key)
		if err != nil {
			return
		}

		allowed, err := s.Limiter.Allow(
			r.Context(),
			hashed_key,
			time.Duration(window_seconds)*time.Second,
			max_reqs,
		)
		if err != nil {
			log.Print(err)
			http.Error(w, "Rate limiter error", http.StatusInternalServerError)
			return
		}
		if !allowed {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) getPlan(
	w http.ResponseWriter,
	r *http.Request,
	hashed_key string,
) (int, int, error) {
	var (
		max_reqs       int
		window_seconds int
	)

	plan_key := fmt.Sprintf(MIDDLEWARE_PLAN_KEY, hashed_key)

	val, err := s.Redis.HGetAll(r.Context(), plan_key).Result()
	if err == nil && len(val) > 0 {
		max_reqs, _ = strconv.Atoi(val["max_reqs"])
		window_seconds, _ = strconv.Atoi(val["window_seconds"])
		return max_reqs, window_seconds, nil
	}

	row, err := s.Queries.GetKeyByHash(r.Context(), hashed_key)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return 0, 0, err
	} else if err != nil {
		log.Print(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return 0, 0, err
	}

	max_reqs = int(row.RequestsPerWindow)
	window_seconds = int(row.WindowSizeSeconds)

	_, _ = s.Redis.HSet(r.Context(), plan_key, map[string]any{
		"max_reqs":       max_reqs,
		"window_seconds": window_seconds,
	}).Result()

	s.Redis.Expire(
		r.Context(),
		plan_key,
		time.Duration(s.Vars.CACHE_TTL_MINUTES)*time.Minute,
	)

	return max_reqs, window_seconds, nil
}
