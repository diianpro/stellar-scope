package main

import (
	"github.com/caarlos0/env/v6"

	"github.com/diianpro/stellar-scope/app"
	"github.com/diianpro/stellar-scope/internal/config"
)

func main() {
	app.New(func() (config.Config, error) {
		cfg := config.Config{}
		if err := env.Parse(&cfg); err != nil {
			return config.Config{}, err
		}
		return cfg, nil
	}).Start()
}
