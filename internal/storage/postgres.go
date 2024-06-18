package storage

import (
	"context"

	"github.com/OurLuv/geograkom/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConn(cfg config.Config) (*pgxpool.Pool, error) {

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
