package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/jollyboss123/finance-tracker/pkg/currency"
	"github.com/jollyboss123/finance-tracker/pkg/expense"
	"github.com/jollyboss123/finance-tracker/pkg/health"
	"github.com/jollyboss123/finance-tracker/pkg/middleware"
	"github.com/jollyboss123/finance-tracker/pkg/rate"
	"github.com/jollyboss123/finance-tracker/pkg/server/response"
	"net/http"
)

func (s *Server) InitDomains() {
	s.initVersion()
	s.initHealth()
	s.initExpense()
}

func (s *Server) initVersion() {
	s.router.Route("/version", func(r chi.Router) {
		r.Use(middleware.Json)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			response.Json(s.l, w, http.StatusOK, map[string]string{
				"version": s.Version,
			})
		})
	})
}

func (s *Server) initHealth() {
	newHealthRepo := health.NewRepo(s.db)
	health.SetupRoutes(s.l, s.router, newHealthRepo)
}

func (s *Server) initExpense() {
	newExpenseRepo := expense.New(s.db)
	newCurrencyRepo := currency.New(s.db)
	newRateRepo := rate.New(s.db)
	er := rate.NewExchangeRates(s.l, newRateRepo)
	expense.SetupRoutes(s.l, s.router, s.validator, newExpenseRepo, newCurrencyRepo, er)
}

func (s *Server) initCurrency() {
	newCurrencyRepo := currency.New(s.db)
	currency.SetupRoutes(s.l, s.router, s.validator, newCurrencyRepo)
}

func (s *Server) initRate() {
	newRateRepo := rate.New(s.db)
	r := rate.NewExchangeRates(s.l, newRateRepo)
	rate.SetupRoutes(s.l, s.router, s.validator, newRateRepo, r, s.cfg)
}
