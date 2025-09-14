package db

import (
	"context"
	"time"

	"github.com/14mdzk/exp/internal/platform/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(cfg.Database.DSN())
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(cfg.Database.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.Database.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.Database.MaxLifetime
	poolConfig.MaxConnIdleTime = time.Minute * time.Duration(cfg.Database.MaxIdleConns)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &PostgresDB{Pool: pool}, nil
}

func (p *PostgresDB) Close() error {
	if p.Pool != nil {
		p.Pool.Close()
	}

	return nil
}

func (p *PostgresDB) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := p.Pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) Stats() *pgxpool.Stat {
	return p.Pool.Stat()
}

// QueryRow executes a query that returns a single row
func (p *PostgresDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return p.Pool.QueryRow(ctx, sql, args...)
}

// Query executes a query that returns multiple rows
func (p *PostgresDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Exec executes a query that doesn't return rows
func (p *PostgresDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	result, err := p.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	return result, nil
}
