package storage

import (
	"log"
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
	conn, err = NewConn(cfg) // fix it
	if err != nil {
		log.Fatalf("failed to init storage: %s", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestCreateRoute(t *testing.T) {
	repo := NewRouteStorage(conn)
	route := model.Route{
		Id:        1,
		Name:      "John",
		Load:      500.4,
		CargoType: "sand",
		IsActual:  true,
	}
	_, err := repo.CreateRoute(route)
	if err != nil {
		t.Error(err)
	}
}

func TestGetRouteById(t *testing.T) {
	repo := NewRouteStorage(conn)
	id := 2
	res, err := repo.GetRouteByID(id)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}
