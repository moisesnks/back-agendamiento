package handlers

import (
	"backend/models"
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func ObtenerPaquetesDestacados(w http.ResponseWriter, r *http.Request) {
	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(`
	SELECT
	PAQ.id AS "ID del Paquete",
    CONCAT(ORI.nombre, ' - ', DEST.nombre) AS "Nombre del Origen y Destino",
    HH.descripcion AS "Que incluye el paquete",
    FP.fechaInit AS "Fecha de Inicio",
    FP.fechaFin AS "Fecha de Fin",
    FP.precioOferta AS Precio
	FROM
    Paquete AS PAQ
	INNER JOIN
    fechaPaquete AS FP ON PAQ.id = FP.id_paquete
	INNER JOIN
    Aeropuerto AS ORI ON PAQ.id_origen = ORI.id
	INNER JOIN
    Aeropuerto AS DEST ON PAQ.id_destino = DEST.id
	INNER JOIN
    habitacionHOTEL AS HH ON PAQ.id_hh = HH.id
	ORDER BY
    RANDOM() 
	LIMIT
    1;
	`)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var paqueteDestacados []models.PaquetesDestacados

	for rows.Next() {
		var paqueteDestacado models.PaquetesDestacados
		err := rows.Scan(&paqueteDestacado.Id, &paqueteDestacado.OrigenDestino, &paqueteDestacado.Detalle, &paqueteDestacado.FechaInicio, &paqueteDestacado.FechaFin, &paqueteDestacado.Precio)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al escanear los resultados", http.StatusInternalServerError)
			return
		}
		fechaInicio, err := time.Parse("2006-01-02T15:04:05Z", paqueteDestacado.FechaInicio)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al parsear la fecha de inicio", http.StatusInternalServerError)
			return
		}

		fechaFin, err := time.Parse("2006-01-02T15:04:05Z", paqueteDestacado.FechaFin)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al parsear la fecha de fin", http.StatusInternalServerError)
			return
		}

		// Formatea las fechas como solo la parte de la fecha (sin la hora)
		paqueteDestacado.FechaInicio = fechaInicio.Format("2006-01-02")
		paqueteDestacado.FechaFin = fechaFin.Format("2006-01-02")

		paqueteDestacados = append(paqueteDestacados, paqueteDestacado)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paqueteDestacados); err != nil {
		log.Fatal(err)
		http.Error(w, "Error al convertir a JSON", http.StatusInternalServerError)
	}
}
