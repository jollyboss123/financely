package server

import (
	"context"
	"crypto/tls"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/jollyboss123/finance-tracker/config"
	"github.com/jollyboss123/finance-tracker/database"
	"github.com/jollyboss123/finance-tracker/pkg/middleware"
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
	cfg     *config.Config

	db *sqlx.DB

	validator *validator.Validate
	cors      *cors.Cors
	tls       *tls.Config
	router    *chi.Mux

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
		log.Printf("Starting API version: %s\n", version)
		opts.Version = version
		return nil
	}
}

func defaultServer() *Server {
	return &Server{
		cfg:    config.New(),
		router: chi.NewRouter(),
	}
}

func (s *Server) Init() {
	s.setCors()
	s.setTls()
	s.NewDatabase()
	s.newValidator()
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
		log.Fatalln("please fill in database credentials in .env file or set in environment variable")
	}

	s.db = database.New(s.cfg.Database)
	s.db.SetMaxOpenConns(s.cfg.Database.MaxConnectionPool)
	s.db.SetMaxIdleConns(s.cfg.Database.MaxIdleConnections)
	s.db.SetConnMaxLifetime(s.cfg.Database.ConnectionMaxLifetime)
}

func (s *Server) newValidator() {
	s.validator = validator.New()
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
	if s.cfg.Api.RequestLog {
		s.router.Use(chiMiddleware.Logger)
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
	log.Printf("Serving at %s:%s\n", s.cfg.Api.Host, s.cfg.Api.Port)
	err := s.httpServer.ListenAndServeTLS(s.cfg.Tls.CertFile, s.cfg.Tls.KeyFile)
	if err != nil {
		log.Printf("Error starting server, %s\n", err)
		os.Exit(1)
	}
}

func gracefulShutdown(ctx context.Context, s *Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down...")

	ctx, shutdown := context.WithTimeout(ctx, s.Config().Api.GracefulTimeout*time.Second)
	defer shutdown()

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
	s.closeResources(ctx)

	return nil
}

func (s *Server) closeResources(ctx context.Context) {
	_ = s.db.Close()
}
