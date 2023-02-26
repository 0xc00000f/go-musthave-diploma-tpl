package storage

import (
	"database/sql"
	"errors"

	"go.uber.org/zap"

	_ "github.com/lib/pq" //nolint:exhaustive

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config/pg"
)

var (
	ErrFailedPgsqlConnectionOpen    = errors.New("failed to open pgsql connection")
	ErrorFailedPgsqlConnectionCheck = errors.New("pgsql ping failed")
)

type Storage struct {
	DB *sql.DB

	Logger *zap.Logger
}

func New(config pg.Config) (*Storage, error) {
	db, err := sql.Open("postgres", config.Dsn)
	if err != nil {
		return nil, ErrFailedPgsqlConnectionOpen
	}

	if err = db.Ping(); err != nil {
		return nil, ErrorFailedPgsqlConnectionCheck
	}

	return &Storage{ //nolint:exhaustruct
		DB: db,
	}, nil
}
