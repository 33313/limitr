package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) setRoutes() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	s.Router.Get("/health", s.handleGetHealth)
	s.Router.Get("/keys", s.handleGetKeys)

	s.Router.Post("/keys", s.handlePostKeys)

	s.Router.Group(func(r chi.Router) {
		r.Use(s.RateLimitMiddleware)

		r.Get("/usage", s.handleGetUsage)

		r.Delete("/keys/{key_id}", s.handleDeleteKey)
	})
}
