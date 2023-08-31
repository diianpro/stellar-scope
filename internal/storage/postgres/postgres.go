package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/url"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, cfg *Config) (*Repository, error) {
	databaseURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	databaseURL.User = url.UserPassword(cfg.Username, cfg.Password)

	connConfig, err := pgxpool.ParseConfig(databaseURL.String())
	if err != nil {
		return nil, err
	}
	connConfig.MinConns = cfg.MinConns
	connConfig.MaxConns = cfg.MaxConns

	conn, err := pgxpool.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	return &Repository{db: conn}, nil

}
