package storage

import (
	"context"

	"github.com/OurLuv/geograkom/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConn(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
