package storage

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrPrepareUsersStatementsFailed = errors.New("failed to prepare users statements")
	ErrUserAlreadyExists            = errors.New("user already exists")
)

const duplicateKeyErrorCode = "23505"

type usersPreparedStatements struct {
	register *sqlx.NamedStmt
}

type Users struct {
	db *sqlx.DB

	prepared usersPreparedStatements
}

func (s *Storage) Users() (*Users, error) {
	u := &Users{
		db:       s.DB,
		prepared: usersPreparedStatements{}, //nolint:exhaustruct
	}

	if err := u.prepareStatements(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *Users) prepareStatements() (err error) {
	u.prepared.register, err = u.db.PrepareNamed(`
		INSERT INTO person (username, password)
		VALUES (:username, :password)
		RETURNING *;
	`)
	if err != nil {
		return errors.Join(ErrPrepareUsersStatementsFailed, err)
	}

	return nil
}

type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

func (u *Users) Register(ctx context.Context, user User) error {
	_, err := u.prepared.register.ExecContext(ctx, user)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok { //nolint:errorlint
			if string(pgErr.Code) == duplicateKeyErrorCode {
				return ErrUserAlreadyExists
			}
		}

		return err //nolint:wrapcheck
	}

	return nil
}
