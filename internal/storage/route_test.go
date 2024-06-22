package storage

import (
	"context"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/OurLuv/geograkom/internal/config"
	"github.com/OurLuv/geograkom/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	cfg := config.Config{StoragePath: "postgres://postgres:admin@localhost:5432/geograkom"}
	conn, err = NewConn(context.Background(), cfg) // fix it
	if err != nil {
		log.Fatalf("failed to init storage: %s", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestCreateRoute(t *testing.T) {
	repo := NewRouteStorage(conn, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	route := model.Route{
		Id:        1,
		Name:      "John",
		Load:      500.4,
		CargoType: "sand",
		IsActual:  true,
	}
	_, err := repo.CreateRoute(context.Background(), route)
	if err != nil {
		t.Error(err)
	}
}

func TestGetRouteById(t *testing.T) {
	repo := NewRouteStorage(conn, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	id := 1
	res, err := repo.GetRouteByID(context.Background(), id)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}
