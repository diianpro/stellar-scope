package postgres

import (
	"context"
	"errors"
	"fmt"
	"path"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, cfg *Config) (*Repository, error) {
	connConfig, err := pgxpool.ParseConfig(cfg.URL)
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

func (r *Repository) Close() {
	r.db.Close()
}

func ApplyMigrate(databaseUrl, migrationsDir string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("could not find migration path")
	}
	dir := path.Join(path.Dir(filename), migrationsDir)

	mig, err := migrate.New(
		fmt.Sprintf("file://%s", dir),
		databaseUrl)
	if err != nil {
		log.Errorf("failed to create migrations instance: %v", err)
		return err
	}

	if err := mig.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Errorf("could not exec migration: %v", err)
		return err
	}
	return nil
}
