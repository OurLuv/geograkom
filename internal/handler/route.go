package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/OurLuv/geograkom/internal/model"
	"github.com/gorilla/mux"
)

// I prefer to describe the interface where I use it
type RouteService interface {
	RegisterRoute(model.Route) (*model.Route, error)
	GetRouteByID(int) (*model.Route, error)
}

// * Register route
func (h *Handler) RegisterRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var route model.Route

	// getting data
	err := json.NewDecoder(r.Body).Decode(&route)
	if err != nil {
		SendError(w, err.Error(), http.StatusBadRequest)
		h.log.Error("bad request", slog.String("err", err.Error()))
		return
	}
	h.log.Debug("getting data", slog.Any("route", route))

	// validating (route_id is not a required field here)
	if route.Name == "" || route.CargoType == "" || route.Load == 0 {
		SendError(w, "did't pass validation", http.StatusBadRequest)
		h.log.Error("did't pass validation", slog.Any("route", route))
		return
	}

	result, err := h.RouteService.RegisterRoute(route)
	if err != nil {
		SendError(w, "server error", http.StatusInternalServerError)
		h.log.Error("server error", slog.String("err", err.Error()))
		return
	}

	// sending response
	var msg string
	if result.SuccesStatusCode == 208 {
		msg = fmt.Sprintf("Route [%d][id] now is not actual anymore", route.Id)
	}
	resp := Response{
		Status:  1,
		Message: msg,
		Data:    result,
	}
	h.log.Debug("Result of registering route", slog.Any("resp", resp))
	w.WriteHeader(result.SuccesStatusCode)
	json.NewEncoder(w).Encode(resp)

}

// * Get request by id
func (h *Handler) GetRouteByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	routeIdStr := mux.Vars(r)["id"]
	routeId, err := strconv.Atoi(routeIdStr)

	// validating
	if err != nil {
		SendError(w, err.Error(), http.StatusBadRequest)
		h.log.Error("bad request", slog.String("err", err.Error()))
		return
	}

	result, err := h.RouteService.GetRouteByID(routeId)
	if err != nil {
		SendError(w, "server error", http.StatusInternalServerError)
		h.log.Error("server error", slog.String("err", err.Error()))
		return
	}

	// route's flag is_actual = FALSE
	if result.SuccesStatusCode == 410 {
		h.log.Debug("Result of registering route", slog.Any("resp", result))
		w.WriteHeader(result.SuccesStatusCode)
		json.NewEncoder(w).Encode(result)
		return
	}

	h.log.Debug("Result of registering route", slog.Any("resp", result))
	w.WriteHeader(result.SuccesStatusCode)
	json.NewEncoder(w).Encode(result)

}
