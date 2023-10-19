package routes

import (
	"backend/api/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func ConfigureRoutes(r *mux.Router) {
	// Habilitar CORS utilizando el paquete rs/cors solo para la ruta /aeropuertos
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://react.lumonidy.studio"},
		AllowedMethods: []string{"GET"},
	})

	// Aplicar el middleware CORS solo a la ruta /aeropuertos
	r.Handle("/aeropuertos", c.Handler(http.HandlerFunc(handlers.ObtenerAeropuertos)))
	r.Handle("/paquetes", c.Handler(http.HandlerFunc(handlers.ObtenerPaquetes)))
	r.Handle("/paquetes/mes", c.Handler(http.HandlerFunc(handlers.ObtenerPaquetesMes)))
	r.Handle("/paquetesoferta", c.Handler(http.HandlerFunc(handlers.ObtenerPaquetesOfertas)))
	r.Handle("/paquetesdestacados", c.Handler(http.HandlerFunc(handlers.ObtenerPaquetesDestacados)))
	r.Handle("/listaeropuertos", c.Handler(http.HandlerFunc(handlers.ListarAeropuertos)))

	// Agrega más configuraciones de rutas aquí si es necesario
}
