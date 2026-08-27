package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func NewPostgresClient(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		pool.Close()
		return nil, fmt.Errorf("set goose dialect: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	if err := goose.Up(db, "migrations"); err != nil {
		db.Close()
		pool.Close()
		return nil, fmt.Errorf("apply migrations: %w", err)
	}

	if err := db.Close(); err != nil {
		pool.Close()
		return nil, fmt.Errorf("close migration db: %w", err)
	}

	return pool, nil
}
