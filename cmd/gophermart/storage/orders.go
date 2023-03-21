package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/0xc00000f/go-musthave-diploma-tpl/lib/libsqlx"
	"github.com/0xc00000f/go-musthave-diploma-tpl/lib/must"
)

type ordersPreparedStatements struct {
	create      *sqlx.NamedStmt
	fetch       *sqlx.NamedStmt
	fetchByUser *sqlx.NamedStmt
}

type Orders struct {
	db *sqlx.DB

	prepared ordersPreparedStatements
}

func (s *Storage) Orders() (*Orders, error) {
	o := &Orders{
		db:       s.DB,
		prepared: ordersPreparedStatements{}, //nolint:exhaustruct
	}

	o.prepareStatements()

	return o, nil
}

func (o *Orders) prepareStatements() {
	o.prepared.create = must.OK(o.db.PrepareNamed(`
		INSERT INTO orders (username, number)
		VALUES (:username, :number)
		RETURNING *;
	`))

	o.prepared.fetch = must.OK(o.db.PrepareNamed(`
		SELECT *
		FROM orders
		WHERE number = ANY(:number)
		LIMIT :limit;
	`))

	o.prepared.fetchByUser = must.OK(o.db.PrepareNamed(`
		SELECT *
		FROM orders
		WHERE username = :username;
	`))
}

type OrderData struct {
	Username    string `db:"username"`
	OrderNumber string `db:"number"`
	Accrual     int64  `db:"accrual"`
	Withdraw    int64  `db:"withdraw"`
	Status      string `db:"status"`
	CreatedTS   int64  `db:"created_ts"`
}

type OrderCreateData struct {
	Username    string `db:"username"`
	OrderNumber string `db:"number"`
}

func (o *Orders) Create(ctx context.Context, data OrderCreateData) (*OrderData, error) {
	rows, err := o.prepared.create.QueryxContext(ctx, data)
	if err != nil {
		return nil, ErrUnexpectedDBError
	}

	result, err := libsqlx.StructScanOneRow[OrderData](rows)
	if err != nil {
		return nil, ErrUnexpectedDBError
	}

	return result, nil
}

type OrderDataMap map[string]*OrderData

func (o *Orders) Fetch(ctx context.Context, numbers []string) (OrderDataMap, error) {
	rows, err := o.prepared.fetch.QueryxContext(
		ctx,
		map[string]any{"number": pq.Array(numbers), "limit": len(numbers)},
	)
	if err != nil {
		return nil, ErrUnexpectedDBError
	}

	result := make(OrderDataMap)

	err = libsqlx.StructScanFn(rows, func(row *OrderData) error {
		result[row.OrderNumber] = row

		return nil
	})

	if err != nil {
		return nil, ErrUnexpectedDBError
	}

	return result, nil
}

func (o *Orders) FetchByUser(ctx context.Context, username string) (OrderDataMap, error) {
	rows, err := o.prepared.fetchByUser.QueryxContext(
		ctx,
		map[string]any{"username": username},
	)
	if err != nil {
		return nil, ErrUnexpectedDBError
	}

	result := make(OrderDataMap)

	err = libsqlx.StructScanFn(rows, func(row *OrderData) error {
		result[row.OrderNumber] = row

		return nil
	})

	if err != nil {
		return nil, ErrUnexpectedDBError
	}

	return result, nil
}
