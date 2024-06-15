package service

import (
	"github.com/OurLuv/geograkom/internal/model"
)

//go:generate mockgen -source=../handler/route.go -destination=mock/mock.go

type RouteRepository interface {
	CreateRoute(route model.Route) (*model.Route, error)
	GetRouteByID(id int) (*model.Route, error)
	DeleteRouteById(id int) error
}

type Route struct {
	repo RouteRepository
}

// * Register route
func (r *Route) RegisterRoute(route model.Route) (*model.Route, error) {
	result, err := r.repo.CreateRoute(route)
	if err != nil {
		return nil, err
	}
	if result.Id != route.Id && route.Id != 0 {
		result.SuccesStatusCode = 208
	} else {
		result.SuccesStatusCode = 200
	}

	return result, nil
}

// * Getting route by id
func (r *Route) GetRouteByID(id int) (*model.Route, error) {

	result, err := r.repo.GetRouteByID(id)
	if err != nil {
		return nil, err
	}
	if !result.IsActual {
		result.SuccesStatusCode = 410
		return result, nil
	}
	result.SuccesStatusCode = 200

	return result, nil
}

// * Deleting routes
func (r *Route) DeleteRoutes(id int) error {
	err := r.repo.DeleteRouteById(id)
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
