package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Api struct {
	Name              string        `default:"financely_api"`
	Host              string        `default:"0.0.0.0"`
	Port              string        `default:"8081"`
	ReadHeaderTimeout time.Duration `split_words:"true" default:"60s"`
	ReadTimeout       time.Duration `split_words:"true" default:"5s"`
	WriteTimeout      time.Duration `split_words:"true" default:"10s"`
	IdleTimeout       time.Duration `split_words:"true" default:"120s"`
	GracefulTimeout   time.Duration `split_words:"true" default:"8s"`

	RequestLog bool `split_words:"true" default:"false"`
}

func API() Api {
	var api Api
	envconfig.MustProcess("API", &api)

	return api
}
