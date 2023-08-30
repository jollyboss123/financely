package health

import (
	"github.com/go-chi/chi/v5"
	"github.com/jollyboss123/finance-tracker/pkg/middleware"
)

func SetupRoutes(router *chi.Mux, repo *repository) *Handler {
	h := NewHandler(repo)

	router.Route("/api/health", func(r chi.Router) {
		r.Use(middleware.Json)

		r.Get("/", h.Health)
		r.Get("/readiness", h.Readiness)
	})

	return h
}
