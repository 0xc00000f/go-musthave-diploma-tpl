package libsqlx

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func StructScanOneRow[T any](rows *sqlx.Rows) (dest *T, err error) {
	defer func() {
		if err2 := rows.Close(); err2 != nil {
			err = errors.Wrap(err2, err.Error())
		}
	}()

	if rows.Next() {
		dest = new(T)

		return dest, rows.StructScan(dest) //nolint:wrapcheck
	}

	return dest, rows.Err() //nolint:wrapcheck
}

type StructScanHandler[T any] func(row *T) error

func StructScanFn[T any](rows *sqlx.Rows, handler StructScanHandler[T]) (err error) {
	defer func() {
		if err2 := rows.Close(); err2 != nil {
			err = errors.Wrap(err2, err.Error())
		}
	}()

	for rows.Next() {
		dest := new(T)

		if err = rows.StructScan(dest); err != nil {
			return
		}

		if err = handler(dest); err != nil {
			return
		}
	}

	return rows.Err() //nolint:wrapcheck
}
