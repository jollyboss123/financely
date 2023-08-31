package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/jollyboss123/finance-tracker/pkg/health"
	"github.com/jollyboss123/finance-tracker/pkg/middleware"
	"github.com/jollyboss123/finance-tracker/pkg/server/response"
	"net/http"
)

func (s *Server) InitDomains() {
	s.initVersion()
	s.initHealth()
}

func (s *Server) initVersion() {
	s.router.Route("/version", func(r chi.Router) {
		r.Use(middleware.Json)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			response.Json(w, http.StatusOK, map[string]string{
				"version": s.Version,
			})
		})
	})
}

func (s *Server) initHealth() {
	newHealthRepo := health.NewRepo(s.db)
	health.SetupRoutes(s.router, newHealthRepo)
}
