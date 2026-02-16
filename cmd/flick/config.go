package main

import (
	"github.com/astrokkidd/flick/pkg/identity"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	JwtSecret            identity.SecretKey `envconfig:"jwt_secret"`
	PostgresUrl          string             `envconfig:"postgres_url"`
	ApiBaseUrl           string             `envconfig:"api_base_address"`
	MessageEncryptionKey string             `envconfig:"message_encryption_key"`
}

func (cfg *Config) Load() {
	envconfig.MustProcess("flick", cfg)
}
