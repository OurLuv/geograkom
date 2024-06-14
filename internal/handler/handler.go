package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type Handler struct {
	RouteService
	log *slog.Logger
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()
	// регистрация(создание) нового маршрута
	r.HandleFunc("/api/route/register", h.RegisterRoute).Methods("POST")

	//(где id - индентификатор маршрута)  - получение данных о маршруте по индентификатору.
	r.HandleFunc("/api/route/{id}", h.GetRouteByID).Methods("GET")

	//удаление маршрутов
	r.HandleFunc("/api/route/", h.DeleteRoutes).Methods("DELETE")

	return r
}

func SendError(w http.ResponseWriter, errorStr string, code int) {
	w.WriteHeader(code)
	response := Response{
		Status: -1,
		Error:  errorStr,
	}
	json.NewEncoder(w).Encode(response)
}

func NewHandler(s RouteService, log *slog.Logger) *Handler {
	return &Handler{
		RouteService: s,
		log:          log,
	}
}
