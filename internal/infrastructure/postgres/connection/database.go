package connection

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Database владеет общим пулом соединений PostgreSQL.
type Database struct {
	Pool *pgxpool.Pool
}

func Open(
	ctx context.Context,
	dsn string,
) (*Database, error) {
	pool, err := openPool(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &Database{
		Pool: pool,
	}, nil
}

func (db *Database) Close() {
	db.Pool.Close()
}
