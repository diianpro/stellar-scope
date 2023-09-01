package config

import (
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"

	"github.com/diianpro/stellar-scope/internal/provider/apod"
	"github.com/diianpro/stellar-scope/internal/storage/postgres"
	"github.com/diianpro/stellar-scope/internal/storage/s3"
)

type Config struct {
	APOD         apod.Config
	Postgres     postgres.Config
	ImageStorage s3.Config
}

// New initialize Config structure
func New() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Errorf("New config: %v", err)
	}
	return &cfg, nil
}
