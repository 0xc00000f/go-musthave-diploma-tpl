package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/structures/status"
	"github.com/0xc00000f/go-musthave-diploma-tpl/lib/libsqlx"
	"github.com/0xc00000f/go-musthave-diploma-tpl/lib/must"
)

type ordersPreparedStatements struct {
	create      *sqlx.NamedStmt
	fetch       *sqlx.NamedStmt
	fetchByUser *sqlx.NamedStmt
	info        *sqlx.NamedStmt
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
		INSERT INTO orders (username, number, status, withdraw)
		VALUES (:username, :number, :status, :withdraw)
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

	o.prepared.info = must.OK(o.db.PrepareNamed(`
		SELECT SUM(accrual) AS balance, SUM(withdraw) as withdraw
		FROM orders
		WHERE username = :username;
	`))
}

type OrderData struct {
	Username    string             `db:"username"`
	OrderNumber string             `db:"number"`
	Accrual     float64            `db:"accrual"`
	Withdraw    float64            `db:"withdraw"`
	Status      status.OrderStatus `db:"status"`
	CreatedTS   int64              `db:"created_ts"`
}

type OrderCreateData struct {
	Username    string             `db:"username"`
	OrderNumber string             `db:"number"`
	Withdraw    float64            `db:"withdraw"`
	Status      status.OrderStatus `db:"status"`
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

type UserInfoData struct {
	Balance  float64 `db:"balance"`
	Withdraw float64 `db:"withdraw"`
}

func (o *Orders) FetchUserInfo(ctx context.Context, username string) (*UserInfoData, error) {
	rows, err := o.prepared.info.QueryxContext(
		ctx,
		map[string]any{"username": username},
	)
	if err != nil {
		return nil, ErrUnexpectedDBError
	}

	result, err := libsqlx.StructScanOneRow[UserInfoData](rows)
	if err != nil {
		return nil, ErrUnexpectedDBError
	}

	return result, nil
}
