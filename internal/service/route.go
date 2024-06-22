package service

import (
	"context"
	"net/http"

	"github.com/OurLuv/geograkom/internal/model"
)

//go:generate mockgen -source=../handler/route.go -destination=mock/mock.go

type RouteRepository interface {
	CreateRoute(ctx context.Context, route model.Route) (*model.Route, error)
	GetRouteByID(ctx context.Context, id int) (*model.Route, error)
	DeleteRouteById(ctx context.Context, id int) error
}

type Route struct {
	repo RouteRepository
}

// * Register route
func (r *Route) RegisterRoute(ctx context.Context, route model.Route) (*model.Route, error) {
	result, err := r.repo.CreateRoute(ctx, route)
	if err != nil {
		return nil, err
	}
	if result.Id != route.Id && route.Id != 0 {
		result.SuccesStatusCode = http.StatusAlreadyReported
	} else {
		result.SuccesStatusCode = http.StatusOK
	}

	return result, nil
}

// * Getting route by id
func (r *Route) GetRouteByID(ctx context.Context, id int) (*model.Route, error) {

	result, err := r.repo.GetRouteByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !result.IsActual {
		result.SuccesStatusCode = http.StatusGone
		return result, nil
	}
	result.SuccesStatusCode = http.StatusOK

	return result, nil
}

// * Deleting routes
func (r *Route) DeleteRoutes(ctx context.Context, id int) error {
	err := r.repo.DeleteRouteById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func NewRouteServcie(repo RouteRepository) *Route {
	return &Route{
		repo: repo,
	}
}
