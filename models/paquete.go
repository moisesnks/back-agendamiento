package models

type Paquete struct {
	ID            int     `json:"id_paquete"`
	Nombre        string  `json:"nombre_paquete"`
	Descripcion   string  `json:"desc_paquete"`
	Detalles      string  `json:"detalle_paquete"`
	FechaInit     string  `json:"fechaInit"`
	FechaFin      string  `json:"fechaFin"`
	CantDias      int     `json:"cant_dias"`
	PrecioViaje   float64 `json:"pr_viaje"`
	PrecioNoche   float64 `json:"pr_noche"`
	PrecioTotal   float64 `json:"pr_total"`
	CiudadOrigen  string  `json:"ciudad_origen"`
	CiudadDestino string  `json:"ciudad_destino"`
	NombreHotel   string  `json:"nombre_hotel"`
	DescripcionHH string  `json:"desc_hh"`
	ServiciosHH   string  `json:"servicios_hh"`
}
