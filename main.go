package main

import (
	"log"
	"net/http"

	"backend/api/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Crear un enrutador utilizando gorilla/mux
	r := mux.NewRouter()

	// Configurar rutas desde routes/routes.go
	routes.ConfigureRoutes(r)

	// Configurar el servidor web para escuchar en el puerto 8080
	port := ":9090"
	log.Printf("Server is ready! Listening on %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
