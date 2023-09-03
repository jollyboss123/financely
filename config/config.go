package config

import (
	"github.com/joho/godotenv"
)

type Config struct {
	Api
	Cors
	Tls

	Cron
	Database
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	return &Config{
		Api:      API(),
		Cors:     NewCors(),
		Tls:      TLS(),
		Cron:     CRON(),
		Database: DataStore(),
	}
}
