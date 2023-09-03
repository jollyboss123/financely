package health

import (
	"github.com/go-chi/chi/v5"
	"github.com/jollyboss123/finance-tracker/pkg/logger"
	"github.com/jollyboss123/finance-tracker/pkg/middleware"
)

func SetupRoutes(l *logger.Logger, router *chi.Mux, repo *repository) *Handler {
	h := NewHandler(l, repo)

	router.Route("/api/health", func(r chi.Router) {
		r.Use(middleware.Json)

		r.Get("/", h.Health)
		r.Get("/readiness", h.Readiness)
	})

	return h
}
