package main

import (
	"context"
	"github.com/jollyboss123/finance-tracker/config"
	"github.com/jollyboss123/finance-tracker/database"
	"github.com/jollyboss123/finance-tracker/pkg/cron"
	"github.com/jollyboss123/finance-tracker/pkg/rate"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.New()
	store := database.New(cfg.Database)
	newRateRepo := rate.New(store)
	r := rate.NewExchangeRates(newRateRepo)
	ctx := context.Background()

	startTime, err := time.Parse("2006-01-02 15:04:05", "2023-09-02 13:30:00")
	if err != nil {
		log.Fatalln(err)
	}
	delay := time.Minute

	jobFunc := func(t time.Time) {
		r.GetRatesRemote(ctx)
	}

	jobID, err := cron.Start("fetch.exchange-rates", startTime, delay, jobFunc)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("started cron job: %s\n", jobID)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down...")

	_, shutdown := context.WithTimeout(ctx, cfg.Api.GracefulTimeout*time.Second)
	defer shutdown()

	store.Close()
}
