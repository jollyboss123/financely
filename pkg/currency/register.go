package currency

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func SetupRoutes(router *chi.Mux, validator *validator.Validate, repo *currencyRepository) *Handler {
	h := NewHandler(repo, validator)

	router.Route("/api/v1/currency", func(r chi.Router) {
		r.Get("/", h.List)
		r.Get("/{currencyID}", h.Get)
		r.Post("/", h.Create)
		r.Put("/{currencyID}", h.Update)
		r.Delete("/{currencyID}", h.Delete)
	})
	return h
}
