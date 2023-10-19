package models

type PaqueteOferta struct {
	OrigenDestino string  `json:"origendestino"`
	Detalle       string  `json:"detalle"`
	FechaInicio   string  `json:"fechainicio"`
	FechaFin      string  `json:"fechafin"`
	Precio        float64 `json:"precio"`
}
