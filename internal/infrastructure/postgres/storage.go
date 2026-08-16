package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Pool *pgxpool.Pool
}

func New(
	ctx context.Context,
	dsn string,
) (*Storage, error) {
	pool, err := connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &Storage{
		Pool: pool,
	}, nil
}

func (s *Storage) Close() {
	s.Pool.Close()
}
