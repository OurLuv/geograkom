package storage

import (
	"context"
	"fmt"

	"github.com/OurLuv/geograkom/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConn(cfg config.Config) (*pgxpool.Pool, error) {

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	var greeting string
	err = pool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		return nil, err
	}
	fmt.Print(greeting)

	return pool, nil
}
