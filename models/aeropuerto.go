package models

type Aeropuerto struct {
	ID         int    `json:"id"`
	Aeropuerto string `json:"aeropuerto"`
	Ciudad     string `json:"ciudad"`
	Pais       string `json:"pais"`
}
