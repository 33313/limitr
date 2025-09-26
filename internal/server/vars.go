package server

import (
	"log"
	"os"
	"strconv"
)

type Vars struct {
	CACHE_TTL_MINUTES           int
	DEFAULT_WINDOW_SIZE_SECONDS int
	DEFAULT_REQUESTS_PER_WINDOW int
}

const (
	FALLBACK_CACHE_TTL_MINUTES   = 60
	FALLBACK_WINDOW_SIZE_SECONDS = 60
	FALLBACK_REQUESTS_PER_WINDOW = 60
)

func getVars() *Vars {
	CACHE_TTL_MINUTES := getEnvIntWithFallback(
		"CACHE_TTL_MINUTES", FALLBACK_CACHE_TTL_MINUTES,
	)
	DEFAULT_WINDOW_SIZE_SECONDS := getEnvIntWithFallback(
		"DEFAULT_WINDOW_SIZE_SECONDS", FALLBACK_WINDOW_SIZE_SECONDS,
	)
	DEFAULT_REQUESTS_PER_WINDOW := getEnvIntWithFallback(
		"DEFAULT_REQUESTS_PER_WINDOW", FALLBACK_REQUESTS_PER_WINDOW,
	)
	return &Vars{
		CACHE_TTL_MINUTES:           CACHE_TTL_MINUTES,
		DEFAULT_WINDOW_SIZE_SECONDS: DEFAULT_WINDOW_SIZE_SECONDS,
		DEFAULT_REQUESTS_PER_WINDOW: DEFAULT_REQUESTS_PER_WINDOW,
	}
}

func getEnvIntWithFallback(name string, fallback int) int {
	res, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		log.Printf("Malformed %s, setting to %v", name, fallback)
		return fallback
	}
	return res
}
