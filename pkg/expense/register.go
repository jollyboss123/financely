package expense

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jollyboss123/finance-tracker/pkg/currency"
	"github.com/jollyboss123/finance-tracker/pkg/logger"
	"github.com/jollyboss123/finance-tracker/pkg/rate"
)

func SetupRoutes(l *logger.Logger, router *chi.Mux, validator *validator.Validate, repo *expenseRepository, cRepo currency.Currency, rates *rate.ExchangeRates) *Handler {
	h := NewHandler(l, repo, cRepo, rates, validator)

	router.Route("/api/v1/expense", func(r chi.Router) {
		r.Get("/", h.List)
		r.Get("/{expenseID}", h.Fetch)
		r.Post("/", h.Create)
		r.Put("/{expenseID}", h.Update)
		r.Delete("/{expenseID}", h.Delete)
		r.Get("/total", h.Total)
		r.Get("/average", h.Average)
	})
	return h
}
