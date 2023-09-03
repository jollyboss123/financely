package expense

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jollyboss123/finance-tracker/pkg/currency"
	"github.com/jollyboss123/finance-tracker/pkg/rate"
)

func SetupRoutes(router *chi.Mux, validator *validator.Validate, repo *expenseRepository, cRepo currency.Currency, rates *rate.ExchangeRates) *Handler {
	h := NewHandler(repo, cRepo, rates, validator)

	router.Route("/api/v1/expense", func(r chi.Router) {
		r.Get("/", h.List)
		r.Get("/{expenseID}", h.Get)
		r.Post("/", h.Create)
		r.Put("/{expenseID}", h.Update)
		r.Delete("/{expenseID}", h.Delete)
		r.Get("/total", h.Total)
	})
	return h
}
