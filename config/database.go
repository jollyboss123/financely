package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Database struct {
	Driver                string        `required:"true"`
	Host                  string        `default:"0.0.0.0"`
	Port                  uint16        `default:"5432"`
	Name                  string        `default:"postgres"`
	TestName              string        `split_words:"true" default:"test"`
	User                  string        `default:"postgres"`
	Pass                  string        `default:"password"`
	SslMode               string        `split_words:"true" default:"disable"`
	MaxConnectionPool     int           `split_words:"true" default:"4"`
	MaxIdleConnections    int           `split_words:"true" default:"4"`
	ConnectionMaxLifetime time.Duration `split_words:"true" default:"300s"`
}

func DataStore() Database {
	var db Database
	envconfig.MustProcess("DB", &db)

	return db
}
