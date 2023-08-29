package config

import "github.com/kelseyhightower/envconfig"

type Api struct {
	Name string `default:"financely_api"`
	Host string `default:"0.0.0.0"`
	Port string `default:"8081"`
}

func API() Api {
	var api Api
	envconfig.MustProcess("API", &api)

	return api
}
