package server

import (
	"context"
	"crypto/tls"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/jollyboss123/finance-tracker/config"
	"github.com/jollyboss123/finance-tracker/database"
	"github.com/jollyboss123/finance-tracker/pkg/logger"
	"github.com/jollyboss123/finance-tracker/pkg/middleware"
	"github.com/jollyboss123/finance-tracker/pkg/postgresstore"
	"github.com/jollyboss123/scs/v2"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Version string
	l       *logger.Logger
	cfg     *config.Config

	db *sqlx.DB

	validator *validator.Validate
	cors      *cors.Cors
	tls       *tls.Config
	router    *chi.Mux

	session       *scs.SessionManager
	sessionCloser *postgresstore.PostgresStore

	httpServer *http.Server
}

type Options func(opts *Server) error

func New(opts ...Options) *Server {
	s := defaultServer()

	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return s
}

func WithVersion(version string) Options {
	return func(opts *Server) error {
		opts.l.Info().Msgf("Starting API version: %s", version)
		opts.Version = version
		return nil
	}
}

func defaultServer() *Server {
	return &Server{
		l:      logger.New(true),
		cfg:    config.New(),
		router: chi.NewRouter(),
	}
}

func (s *Server) Init() {
	s.setCors()
	s.setTls()
	s.NewDatabase()
	s.newValidator()
	s.newAuthentication()
	s.newRouter()
	s.setGlobalMiddleware()
	s.InitDomains()
}

func (s *Server) setCors() {
	s.cors = cors.New(
		cors.Options{
			AllowedOrigins: s.cfg.Cors.AllowedOrigins,
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		})
}

func (s *Server) setTls() {
	s.tls = &tls.Config{
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},

		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}

func (s *Server) NewDatabase() {
	if s.cfg.Database.Driver == "" {
		s.l.Error().Msg("please fill in database credentials in .env file or set in environment variable")
	}

	s.db = database.New(s.l, s.cfg.Database)
	s.db.SetMaxOpenConns(s.cfg.Database.MaxConnectionPool)
	s.db.SetMaxIdleConns(s.cfg.Database.MaxIdleConnections)
	s.db.SetConnMaxLifetime(s.cfg.Database.ConnectionMaxLifetime)
}

func (s *Server) newValidator() {
	s.validator = validator.New()
}

func (s *Server) newAuthentication() {
	manager := scs.New()
	manager.Store = postgresstore.New(s.db.DB, s.l)
	manager.CtxStore = postgresstore.New(s.db.DB, s.l)
	manager.Lifetime = s.cfg.Session.Duration
	manager.Cookie.Name = s.cfg.Session.Name
	manager.Cookie.Domain = s.cfg.Session.Domain
	manager.Cookie.HttpOnly = s.cfg.Session.HttpOnly
	manager.Cookie.Path = s.cfg.Session.Path
	manager.Cookie.Persist = true
	manager.Cookie.SameSite = http.SameSite(s.cfg.Session.SameSite)
	manager.Cookie.Secure = s.cfg.Session.Secure

	s.sessionCloser = postgresstore.NewWithCleanupInterval(s.db.DB, s.l, 30*time.Minute)

	s.session = manager
}

func (s *Server) newRouter() {
	s.router = chi.NewRouter()
}

func (s *Server) setGlobalMiddleware() {
	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message": "endpoint not found"}`))
	})
	s.router.Use(s.cors.Handler)
	s.router.Use(middleware.Json)
	s.router.Use(middleware.LoadAndSave(s.session))
	if s.cfg.Api.RequestLog {
		s.router.Use(middleware.RequestLog(s.l, s.session))
	}
}

func (s *Server) Run() {
	s.httpServer = &http.Server{
		Addr:              s.cfg.Api.Host + ":" + s.cfg.Api.Port,
		Handler:           s.router,
		ReadHeaderTimeout: s.cfg.Api.ReadHeaderTimeout,
		ReadTimeout:       s.cfg.Api.ReadTimeout,
		WriteTimeout:      s.cfg.Api.WriteTimeout,
		IdleTimeout:       s.cfg.Api.IdleTimeout,
		TLSConfig:         s.tls,
	}

	go func() {
		start(s)
	}()

	_ = gracefulShutdown(context.Background(), s)
}

func (s *Server) Config() *config.Config {
	return s.cfg
}

func start(s *Server) {
	s.l.Info().Msgf("Serving at %s:%s", s.cfg.Api.Host, s.cfg.Api.Port)

	err := s.httpServer.ListenAndServeTLS(s.cfg.Tls.CertFile, s.cfg.Tls.KeyFile)
	if err != nil {
		s.l.Error().Err(err).Msg("Error starting server")
		os.Exit(1)
	}
}

func gracefulShutdown(ctx context.Context, s *Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	s.l.Info().Msg("Shutting down...")

	ctx, shutdown := context.WithTimeout(ctx, s.Config().Api.GracefulTimeout*time.Second)
	defer shutdown()

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		s.l.Err(err)
	}
	s.closeResources()

	return nil
}

func (s *Server) closeResources() {
	_ = s.db.Close()
	s.sessionCloser.StopCleanup()
}
