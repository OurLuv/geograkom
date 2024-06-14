package storage

import (
	"context"

	"github.com/OurLuv/geograkom/internal/model"
	"github.com/jackc/pgx/v5"
)

type RouteStorage struct {
	conn *pgx.Conn
}

func (rs *RouteStorage) CreateRoute(route model.Route) (*model.Route, error) {
	ctx := context.Background()
	tx, err := rs.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var count int
	query := "SELECT COUNT(*) FROM routes WHERE route_id = $1"
	err = tx.QueryRow(ctx, query, route.Id).Scan(&count)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		query = "INSERT INTO routes (route_id, route_name, load, cargo_type, is_actual)" +
			"VALUES ($1, $2, $3, $4, $5);"
		if _, err = tx.Exec(ctx, query, route.Id, route.Name, route.Load, route.CargoType, route.IsActual); err != nil {
			return nil, err
		}
	} else {
		query = "UPDATE routes SET is_actual = false WHERE route_id = $1;"
		if _, err = tx.Exec(ctx, query, route.Id); err != nil {
			return nil, err
		}
		query = "INSERT INTO routes (route_name, load, cargo_type, is_actual)" +
			"VALUES ($1, $2, $3, $4) RETURNING route_id;"
		if err = tx.QueryRow(ctx, query, route.Name, route.Load, route.CargoType, route.IsActual).Scan(&route.Id); err != nil {
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func NewRouteStorage(conn *pgx.Conn) *RouteStorage {
	return &RouteStorage{
		conn: conn,
	}
}
