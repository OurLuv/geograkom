package storage

import (
	"context"
	"fmt"

	"github.com/OurLuv/geograkom/internal/config"
	"github.com/jackc/pgx/v5"
)

func NewConn(cfg config.Config) (*pgx.Conn, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		return nil, err
	}
	fmt.Print(greeting)

	return conn, nil
}
