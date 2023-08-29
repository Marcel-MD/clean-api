package config

import (
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Host        string `env:"HOST" envDefault:"localhost"`
	Port        string `env:"SERVER_PORT" envDefault:"8080"`
	AllowOrigin string `env:"ALLOW_ORIGIN" envDefault:"*"`
	Env         string `env:"ENV" envDefault:"dev"`

	AccessTokenSecret    string        `env:"ACCESS_TOKEN_SECRET" envDefault:"SecretAccessSecretAccess"`
	AccessTokenLifespan  time.Duration `env:"ACCESS_TOKEN_LIFESPAN" envDefault:"1h"`
	RefreshTokenSecret   string        `env:"REFRESH_TOKEN_SECRET" envDefault:"SecretRefreshSecretRefresh"`
	RefreshTokenLifespan time.Duration `env:"REFRESH_TOKEN_LIFESPAN" envDefault:"72h"`
	DatabaseUrl          string        `env:"DATABASE_URL" envDefault:"postgres://postgres:password@postgres:5432/clean-api"`
}

func NewConfig() (Config, error) {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load .env file.")
	}

	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
