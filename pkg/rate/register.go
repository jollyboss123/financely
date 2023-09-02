package rate

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func SetupRoutes(router *chi.Mux, validator *validator.Validate, repo *rateRepository) *Handler {
	h := NewHandler(repo, validator)

	router.Route("/api/v1/rate", func(r chi.Router) {
		r.Put("/reschedule", h.Reschedule)
	})

	return h
}
