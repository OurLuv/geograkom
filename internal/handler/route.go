package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/OurLuv/geograkom/internal/model"
	"github.com/gorilla/mux"
)

// I prefer to describe the interface where I use it
type RouteService interface {
	RegisterRoute(model.Route) (*model.Route, error)
	GetRouteByID(int) (*model.Route, error)
	DeleteRoutes(id int) error
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
	if result.SuccesStatusCode == http.StatusAlreadyReported {
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
	if result.SuccesStatusCode == http.StatusGone {
		h.log.Debug("Result of registering route", slog.Any("resp", result))
		w.WriteHeader(result.SuccesStatusCode)
		resp := Response{
			Status:  2,
			Message: "Route is not actual",
			Error:   "",
			Data:    nil,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	h.log.Debug("Result of getting route", slog.Any("resp", result))
	w.WriteHeader(result.SuccesStatusCode)
	json.NewEncoder(w).Encode(result)

}

// * Deleting routes
func (h *Handler) DeleteRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var routeIds []int

	// getting id's
	if err := json.NewDecoder(r.Body).Decode(&routeIds); err != nil {
		SendError(w, err.Error(), http.StatusBadRequest)
		h.log.Error("bad request", slog.String("err", err.Error()))
		return
	}

	ch := make(chan int)
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	go func() {
		for id := range ch {
			wg.Add(1)
			id := id
			go func() {
				defer wg.Done()
				h.log.Debug("passing to service", slog.Any("id", id))
				err := h.RouteService.DeleteRoutes(id)
				if err != nil {
					h.log.Error("error from DB", slog.String("err", err.Error()))
					return
				}
				h.log.Debug("deleted from db", slog.Any("id", id))
			}()
		}
	}()
	// [1, 3, 4]
	for _, v := range routeIds {
		h.log.Debug("added to chan", slog.Any("id", v))
		ch <- v
	}
	close(ch)
	resp := Response{
		Status:  1,
		Message: "удаление маршрутов принято в обработку",
		Error:   "",
		Data:    nil,
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}
