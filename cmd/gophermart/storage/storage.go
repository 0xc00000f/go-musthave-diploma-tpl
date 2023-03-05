package storage

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	_ "github.com/lib/pq" //revive:disable-line:blank-imports

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config/pg"
)

var (
	ErrFailedPgsqlConnectionOpen    = errors.New("failed to open pgsql connection")
	ErrorFailedPgsqlConnectionCheck = errors.New("pgsql ping failed")
)

type Storage struct {
	DB *sqlx.DB

	Logger *zap.Logger
}

func New(config pg.Config) (*Storage, error) {
	db, err := sqlx.Connect("postgres", config.Dsn)
	if err != nil {
		return nil, ErrFailedPgsqlConnectionOpen
	}

	if err = db.Ping(); err != nil {
		return nil, ErrorFailedPgsqlConnectionCheck
	}

	if err = createUserTable(db); err != nil {
		return nil, fmt.Errorf("failed to create user table: %w", err)
	}

	return &Storage{ //nolint:exhaustruct
		DB: db,
	}, nil
}

func createUserTable(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS person (
			username text PRIMARY KEY,
			password text
		)
	`

	if _, err := db.Exec(query); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}
