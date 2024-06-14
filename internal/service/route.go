package service

import (
	"github.com/OurLuv/geograkom/internal/model"
)

type RouteRepository interface {
	CreateRoute(route model.Route) (*model.Route, error)
	GetRouteByID(id int) (*model.Route, error)
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

func NewRouteServcie(repo RouteRepository) *Route {
	return &Route{
		repo: repo,
	}
}
