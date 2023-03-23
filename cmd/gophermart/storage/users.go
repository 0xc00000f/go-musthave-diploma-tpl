package storage

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/0xc00000f/go-musthave-diploma-tpl/lib/libsqlx"
	"github.com/0xc00000f/go-musthave-diploma-tpl/lib/must"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUnexpectedDBError = errors.New("unexpected db error")
)

const duplicateKeyErrorCode = "23505"

type usersPreparedStatements struct {
	register *sqlx.NamedStmt
	fetch    *sqlx.NamedStmt
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

	u.prepareStatements()

	return u, nil
}

func (u *Users) prepareStatements() {
	u.prepared.register = must.OK(u.db.PrepareNamed(`
		INSERT INTO person (username, password)
		VALUES (:username, :password)
		RETURNING *;
	`))

	u.prepared.fetch = must.OK(u.db.PrepareNamed(`
		SELECT *
		FROM person
		WHERE username = ANY(:username)
		LIMIT :limit;
	`))
}

type UserData struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

func (u *Users) Register(ctx context.Context, user UserData) error {
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

type UserDataMap map[string]*UserData

func (u *Users) Fetch(ctx context.Context, usernames []string) (UserDataMap, error) {
	rows, err := u.prepared.fetch.QueryxContext(
		ctx,
		map[string]any{"username": pq.Array(usernames), "limit": len(usernames)},
	)
	if err != nil {
		return nil, ErrUnexpectedDBError
	}

	result := make(UserDataMap)

	err = libsqlx.StructScanFn(rows, func(row *UserData) error {
		result[row.Username] = row

		return nil
	})

	if err != nil {
		return nil, ErrUnexpectedDBError
	}

	return result, nil
}
