package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"log/slog"

	"github.com/OurLuv/geograkom/internal/handler"
	"github.com/OurLuv/geograkom/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRouteService struct {
	mock.Mock
}

// DeleteRoutes implements handler.RouteService.
func (m *MockRouteService) DeleteRoutes(id int) error {
	panic("unimplemented")
}

func (m *MockRouteService) RegisterRoute(route model.Route) (*model.Route, error) {
	args := m.Called(route)
	return args.Get(0).(*model.Route), args.Error(1)
}

func (m *MockRouteService) GetRouteByID(id int) (*model.Route, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Route), args.Error(1)
}

func TestRegisterRoute(t *testing.T) {
	mockS := new(MockRouteService)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	h := handler.NewHandler(mockS, logger)

	route := model.Route{
		Name:             "Test Route",
		CargoType:        "Type1",
		Load:             23.23,
		SuccesStatusCode: http.StatusOK,
	}

	mockS.On("RegisterRoute", route).Return(&route, nil)

	router := h.InitRoutes()
	requestBody, _ := json.Marshal(route)
	req, err := http.NewRequest("POST", "/api/route/register", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response handler.Response
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	// check data type & assert
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)

	// actual model route from response.Data
	actualRoute := model.Route{
		Id:               int(data["id"].(float64)),
		Name:             data["name"].(string),
		CargoType:        data["cargo_type"].(string),
		Load:             data["load"].(float64),
		IsActual:         data["is_actual"].(bool),
		SuccesStatusCode: http.StatusOK,
	}

	assert.Equal(t, route, actualRoute)
	mockS.AssertExpectations(t)
}

func TestGetRouteByID(t *testing.T) {
	mockService := new(MockRouteService)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	h := handler.NewHandler(mockService, logger)

	expectedRoute := model.Route{
		Id:               1,
		Name:             "Test Route",
		CargoType:        "Type1",
		Load:             100,
		SuccesStatusCode: http.StatusOK,
	}

	mockService.On("GetRouteByID", 1).Return(&expectedRoute, nil)

	router := h.InitRoutes()
	req, err := http.NewRequest("GET", "/api/route/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response model.Route
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, expectedRoute, response)
	mockService.AssertExpectations(t)
}
