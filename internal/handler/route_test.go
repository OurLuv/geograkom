package handler_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/OurLuv/geograkom/internal/handler"
	"github.com/OurLuv/geograkom/internal/model"
	mock_service "github.com/OurLuv/geograkom/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetRouteById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockRouteService, ctx context.Context, id int)

	testCases := []struct {
		name               string
		routeId            int
		expectedStatusCode int
		mockBehavior       mockBehavior
	}{
		{
			name:               "OK",
			routeId:            1,
			expectedStatusCode: http.StatusOK,
			mockBehavior: func(s *mock_service.MockRouteService, ctx context.Context, id int) {
				s.EXPECT().GetRouteByID(ctx, id).Return(&model.Route{
					Id:               1,
					Name:             "James",
					Load:             443.4,
					CargoType:        "sand",
					IsActual:         true,
					SuccesStatusCode: http.StatusOK,
				}, nil)
			},
		},
		{
			name:               "Is Not Actual",
			routeId:            1,
			expectedStatusCode: http.StatusGone,
			mockBehavior: func(s *mock_service.MockRouteService, ctx context.Context, id int) {
				s.EXPECT().GetRouteByID(ctx, id).Return(&model.Route{
					Id:               1,
					Name:             "James",
					Load:             443.4,
					CargoType:        "sand",
					IsActual:         false,
					SuccesStatusCode: http.StatusGone,
				}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			s := mock_service.NewMockRouteService(c)
			tc.mockBehavior(s, context.Background(), tc.routeId)
			h := handler.NewHandler(s, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			})))

			//Test server
			r := mux.NewRouter()
			r.HandleFunc("/api/route/{id}", h.GetRouteByID).Methods("GET")

			// Test request
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/route/%d", tc.routeId)
			req := httptest.NewRequest("GET", url, nil)

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}
