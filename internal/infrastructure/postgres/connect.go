package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, url)
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, err
}
