package db_adapter

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Query interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type Transaction interface {
	Query
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type withTx[T any] interface {
	WithTx(tx Transaction) T
}
