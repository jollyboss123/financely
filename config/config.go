package config

import (
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Api
	Cors
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
		Cors:     NewCors(),
		Tls:      TLS(),
		Database: DataStore(),
	}
}
