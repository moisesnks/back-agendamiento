package handlers

import (
	"backend/api/models"
	"backend/api/utils"
	"encoding/json"
	"log"
	"net/http"
)

func ListarAeropuertos(w http.ResponseWriter, r *http.Request) {
	// Establece la conexión con la base de datos PostgreSQL utilizando la URL de conexión
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
		Aeropuerto.id,
		CONCAT(Aeropuerto.nombre, ', ', Ciudad.nombre, ', ', Pais.nombre) as Aeropuerto_Origen
	FROM
		Aeropuerto
	JOIN
		Ciudad ON Aeropuerto.ciudad_id = Ciudad.id
	JOIN
		Pais ON Ciudad.pais_id = Pais.id;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Slice para almacenar los aeropuertos
	var aeropuertos []models.Aeropuertos

	// Itera a través de los resultados y agrega a la slice
	for rows.Next() {
		var aeropuerto models.Aeropuertos
		err := rows.Scan(&aeropuerto.ID, &aeropuerto.Aeropuerto)
		if err != nil {
			log.Fatal(err)
		}
		aeropuertos = append(aeropuertos, aeropuerto)
	}

	// Convierte los resultados a JSON y envía la respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(aeropuertos)
}
