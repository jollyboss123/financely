package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Cron struct {
	ExchangeRatesJobName string        `split_words:"true" default:"fetch.exchange-rates"`
	ExchangeRatesEnabled bool          `split_words:"true" default:"true"`
	ExchangeRatesStart   string        `split_words:"true"`
	ExchangeRatesDelay   time.Duration `split_words:"true" default:"24h"`
}

func CRON() Cron {
	var cron Cron
	envconfig.MustProcess("CRON", &cron)

	return cron
}
