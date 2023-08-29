package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jollyboss123/finance-tracker/config"
	"github.com/jollyboss123/finance-tracker/database"
	"github.com/jollyboss123/finance-tracker/pkg/expenses"
	"github.com/jollyboss123/finance-tracker/pkg/server"
	"log"
	"os"
	"os/signal"
	"time"
)

//var (
//	CertFile    = env.Get("CERT_FILE")
//	KeyFile     = env.Get("KEY_FILE")
//	ServiceAddr = env.Get("SERVICE_ADDR")
//)

func main() {
	l := log.New(os.Stdout, "Financely ", log.LstdFlags)
	cfg := config.New()
	api := cfg.Api
	tls := cfg.Tls
	database.New(cfg.Database)

	e := expenses.NewExpense(l)

	sm := mux.NewRouter()
	e.SetupRoutes(sm)

	srv := server.New(sm, ":"+api.Port)

	go func() {
		l.Printf("server starting on: %s\n", api.Port)
		err := srv.ListenAndServeTLS(tls.CertFile, tls.KeyFile)
		if err != nil {
			l.Printf("Error starting server, %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	l.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		l.Printf("HTTP server Shutdown: %v", err)
	} else {
		l.Println("HTTP server shut down gracefully.")
	}
}
