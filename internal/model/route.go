package model

type Route struct {
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	Load             float64 `json:"load"`
	CargoType        string  `json:"cargo_type"`
	IsActual         bool    `json:"is_actual"`
	SuccesStatusCode int
}
