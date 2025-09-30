package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rate-limiter/internal/db"
	"rate-limiter/internal/keys"
	"rate-limiter/internal/limiter"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ResponseApiKey struct {
	ApiKey string `json:"api_key"`
}

type ResponseUsage struct {
	Limit     int `json:"limit"`
	Used      int `json:"used"`
	Remaining int `json:"remaining"`
}

func (s *Server) handlePostKeys(w http.ResponseWriter, r *http.Request) {
	key, err := keys.GenerateKey()
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to generate key", http.StatusInternalServerError)
		return
	}
	var (
		window_size_seconds int
		requests_per_window int
	)
	if res, err := strconv.Atoi(r.Header.Get("x-window-size-seconds")); err != nil {
		window_size_seconds = s.Vars.DEFAULT_WINDOW_SIZE_SECONDS
	} else {
		window_size_seconds = res
	}
	if res, err := strconv.Atoi(r.Header.Get("x-requests-per-window")); err != nil {
		requests_per_window = s.Vars.DEFAULT_REQUESTS_PER_WINDOW
	} else {
		requests_per_window = res
	}

	_, err = s.Queries.CreateKey(r.Context(), db.CreateKeyParams{
		HashedKey:         keys.HashKey(key),
		WindowSizeSeconds: int32(window_size_seconds),
		RequestsPerWindow: int32(requests_per_window),
	})
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to create key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := ResponseApiKey{
		ApiKey: key,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleGetHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.DB.Ping(r.Context()); err != nil {
		log.Print(err)
		http.Error(w, "Failed to ping database", http.StatusServiceUnavailable)
		return
	}
	if err := s.Redis.Ping(r.Context()).Err(); err != nil {
		log.Print(err)
		http.Error(w, "Failed to ping Redis", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetKeys(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Queries.ListKeys(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to list keys", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rows)
}

func (s *Server) handleDeleteKey(w http.ResponseWriter, r *http.Request) {
	key_uuid_string := chi.URLParam(r, "key_id")
	var key_uuid pgtype.UUID
	if err := key_uuid.Scan(key_uuid_string); err != nil {
		log.Print(err)
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	if err := s.Queries.DeleteKey(r.Context(), key_uuid); err != nil {
		log.Print(err)
		http.Error(w, "Failed to delete key", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleGetUsage(w http.ResponseWriter, r *http.Request) {
	api_key := r.Header.Get("x-api-key")
	if api_key == "" {
		http.Error(w, "No API key provided", http.StatusUnauthorized)
		return
	}
	hashed_key := keys.HashKey(api_key)

	max_reqs, window_seconds, err := s.getPlan(w, r, hashed_key)
	if err != nil {
		return
	}

	now := time.Now().UnixMilli()
	window_start_ms := now - (int64(window_seconds) * 1000)
	count, err := s.Redis.ZCount(
		r.Context(),
		fmt.Sprintf(limiter.LIMITER_CACHE_KEY, hashed_key),
		strconv.FormatInt(window_start_ms, 10),
		strconv.FormatInt(now, 10),
	).Result()

	w.Header().Set("Content-Type", "application/json")
	data := ResponseUsage{
		Limit:     max_reqs,
		Used:      int(count),
		Remaining: max_reqs - int(count),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
