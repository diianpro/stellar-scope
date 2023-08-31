package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/diianpro/stellar-scope/internal/provider/apod"
	"github.com/diianpro/stellar-scope/internal/storage/postgres"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	APOD     *apod.Config
	Postgres postgres.Config
}

// New initialize Config structure
func New() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Errorf("New config: %v", err)
	}
	return &cfg, nil
}
