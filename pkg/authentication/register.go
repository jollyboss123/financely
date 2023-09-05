package authentication

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jollyboss123/finance-tracker/pkg/logger"
	"github.com/jollyboss123/finance-tracker/pkg/middleware"
	"github.com/jollyboss123/scs/v2"
)

func SetupRoutes(l *logger.Logger, router *chi.Mux, validator *validator.Validate, session *scs.SessionManager, repo User) {
	h := NewHandler(l, validator, repo, session)

	router.Post("/api/v1/login", h.Login)
	router.Post("/api/v1/register", h.Register)

	router.Route("/api/v1/restricted", func(r chi.Router) {
		r.Use(middleware.Authenticate(session))
		r.Get("/csrf", h.Csrf)
		r.Get("/", h.Protected)
		r.Get("/me", h.Me)
	})
}
