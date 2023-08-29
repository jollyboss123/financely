package config

import "github.com/kelseyhightower/envconfig"

type Tls struct {
	//ServiceAddr string `split_words:"true" required:"true"`
	CertFile string `required:"true" split_words:"true"`
	KeyFile  string `required:"true" split_words:"true"`
}

func TLS() Tls {
	var tls Tls
	envconfig.MustProcess("TLS", &tls)

	return tls
}
