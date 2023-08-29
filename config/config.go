package config

import (
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Api
	Tls

	Database
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	return &Config{
		Api:      API(),
		Tls:      TLS(),
		Database: DataStore(),
	}
}
