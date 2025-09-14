package db_adapter

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transactor[T any] interface {
	InTx(ctx context.Context, fn func(T) error) error
}

type impl[T any] struct {
	db *pgxpool.Pool
	wt withTx[T]
}

type Trm[T any] interface {
	InTx(ctx context.Context, fn func(T) error) error
}

func NewTrm[T withTx[T]](db *pgxpool.Pool, wt T) *impl[T] {
	return &impl[T]{db, wt}
}

func (t *impl[T]) InTx(ctx context.Context, fn func(repo T) error) error {
	tx, err := t.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	repo := t.wt.WithTx(tx)
	if err := fn(repo); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
