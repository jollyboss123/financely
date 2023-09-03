package main

import (
	"context"
	"github.com/jollyboss123/finance-tracker/config"
	"github.com/jollyboss123/finance-tracker/database"
	"github.com/jollyboss123/finance-tracker/pkg/cron"
	"github.com/jollyboss123/finance-tracker/pkg/logger"
	"github.com/jollyboss123/finance-tracker/pkg/rate"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.New()
	l := logger.New(true)
	store := database.New(l, cfg.Database)
	newRateRepo := rate.New(store)
	r := rate.NewExchangeRates(newRateRepo)
	ctx := context.Background()

	startTime, err := time.Parse("2006-01-02 15:04:05", cfg.Cron.ExchangeRatesStart)
	if err != nil {
		l.Fatal().Err(err)
	}
	delay := cfg.Cron.ExchangeRatesDelay

	jobFunc := func(t time.Time) {
		r.GetRatesRemote(ctx)
	}

	jobID, err := cron.Start(l, cfg.Cron.ExchangeRatesJobID, startTime, delay, jobFunc)
	if err != nil {
		l.Fatal().Err(err)
	}
	l.Info().Str("id", jobID).Msgf("started at: %v with delay: %v", jobID, startTime, delay)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	l.Info().Msg("Shutting down...")

	_, shutdown := context.WithTimeout(ctx, cfg.Api.GracefulTimeout*time.Second)
	defer shutdown()

	store.Close()
}
