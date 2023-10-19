package handlers

import (
	"backend/api/models"
	"backend/api/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func ObtenerAeropuertos(w http.ResponseWriter, r *http.Request) {
	codigoAeroStr := r.URL.Query().Get("id_aeropuerto")
	nombreAeropuerto := r.URL.Query().Get("nombre_aeropuerto")
	nombrePais := r.URL.Query().Get("nombre_pais")
	nombreCiudad := r.URL.Query().Get("nombre_ciudad")

	codigoAero, err := strconv.Atoi(codigoAeroStr)
	if err != nil {
		http.Error(w, "El parámetro 'id' debe ser un número válido", http.StatusBadRequest)
		return
	}

	// Abre la conexión con la base dse datos
	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Consulta SQL para obtener aeropuertos
	rows, err := db.Query(`
    SELECT
    A.id AS "Código Aeropuerto",
	A.nombre AS "Nombre Aeropuerto",
	C.nombre AS "Ciudad",
    P.nombre AS "País"
FROM
    Aeropuerto AS A
INNER JOIN
    Ciudad AS C ON A.ciudad_id = C.id
INNER JOIN
    Pais AS P ON C.pais_id = P.id
WHERE
    ($1 = 0 OR A.id = $1)
    AND ($2 = '' OR A.nombre = $2)
    AND ($3 = '' OR P.nombre = $3)
    AND ($4 = '' OR C.nombre = $4);
`, codigoAero, nombreAeropuerto, nombrePais, nombreCiudad)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Slice para almacenar los aeropuertos
	var aeropuertos []models.Aeropuerto

	// Itera a través de los resultados y agrega a la slice
	for rows.Next() {
		var aeropuerto models.Aeropuerto
		err := rows.Scan(&aeropuerto.ID, &aeropuerto.Aeropuerto, &aeropuerto.Ciudad, &aeropuerto.Pais)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al escanear resultados", http.StatusInternalServerError)
			return
		}
		aeropuertos = append(aeropuertos, aeropuerto)
	}

	// Convierte los resultados a JSON y envía la respuesta
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(aeropuertos); err != nil {
		log.Fatal(err)
		http.Error(w, "Error al convertir a JSON", http.StatusInternalServerError)
	}
}
