package expense

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func SetupRoutes(router *chi.Mux, validator *validator.Validate, repo *expenseRepository) *Handler {
	h := NewHandler(repo, validator)

	router.Route("/api/v1/expense", func(r chi.Router) {
		r.Get("/", h.List)
		r.Get("/{expenseID}", h.Get)
		r.Post("/", h.Create)
		r.Put("/{expenseID}", h.Update)
		r.Delete("/{expenseID}", h.Delete)
	})
	return h
}
