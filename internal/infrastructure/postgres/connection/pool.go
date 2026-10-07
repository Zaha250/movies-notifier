package connection

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func openPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("создание пула PostgreSQL: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("проверка подключения к PostgreSQL: %w", err)
	}

	return pool, nil
}
