package main

import (
	"context"
	"github.com/jollyboss123/finance-tracker/config"
	"github.com/jollyboss123/finance-tracker/database"
	"github.com/jollyboss123/finance-tracker/pkg/cron"
	"github.com/jollyboss123/finance-tracker/pkg/rate"
	"log"
	"time"
)

func main() {
	cfg := config.New()
	store := database.New(cfg.Database)
	newRateRepo := rate.New(store)
	r := rate.NewExchangeRates(newRateRepo)

	startTime, err := time.Parse("2006-01-02 15:04:05", "2023-09-02 13:30:00")
	if err != nil {
		panic(err)
	}
	delay := time.Minute

	for t := range cron.Cron(context.Background(), startTime, delay) {
		r.GetRatesRemote(context.Background())
		log.Printf("cron run at %v", t.Format("2006-01-02 15:04:05"))
	}
}
