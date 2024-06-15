package storage

import (
	"context"
	"log/slog"

	"github.com/OurLuv/geograkom/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RouteStorage struct {
	conn *pgxpool.Pool
	log  *slog.Logger
}

// * Creating route
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
	rs.log.Debug("Recieved data", slog.Any("route", route))
	if count == 0 {
		if route.Id == 0 {
			query = "INSERT INTO routes (route_name, load, cargo_type, is_actual)" +
				"VALUES ($1, $2, $3, $4) RETURNING route_id;"
			if err = tx.QueryRow(ctx, query, route.Name, route.Load, route.CargoType, route.IsActual).Scan(&route.Id); err != nil {
				return nil, err
			}
		} else {
			query = "INSERT INTO routes (route_id, route_name, load, cargo_type, is_actual)" +
				"VALUES ($1, $2, $3, $4, $5);"
			if _, err = tx.Exec(ctx, query, route.Id, route.Name, route.Load, route.CargoType, route.IsActual); err != nil {
				return nil, err
			}
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

// * Getting route by id
func (rs *RouteStorage) GetRouteByID(id int) (*model.Route, error) {
	query := "SELECT * FROM routes WHERE route_id = $1"
	var result model.Route
	err := rs.conn.QueryRow(context.Background(), query, id).Scan(&result.Id, &result.Name, &result.Load, &result.CargoType, &result.IsActual)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// * Deleting route by id
func (rs *RouteStorage) DeleteRouteById(id int) error {
	query := "DELETE FROM routes WHERE route_id = $1"
	_, err := rs.conn.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}

func NewRouteStorage(conn *pgxpool.Pool, log *slog.Logger) *RouteStorage {
	return &RouteStorage{
		conn: conn,
		log:  log,
	}
}
