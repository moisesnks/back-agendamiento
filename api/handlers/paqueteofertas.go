package handlers

import (
	"backend/api/models"
	"backend/api/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func ObtenerPaquetesOfertas(w http.ResponseWriter, r *http.Request) {
	//CONE CON LA BASE DE DATOS
	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(`
    SELECT
    CONCAT(ORI.nombre, ' - ', DEST.nombre) AS "Nombre del Origen y Destino", 
    HH.descripcion AS "Que incluye el paquete",
    FP.fechaInit AS "Fecha de Inicio",
    FP.fechaFin AS "Fecha de Fin",
    FP.precioOferta AS "Precio"
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
    WHERE
        FP.precioOferta != 0
    `)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var paquetesOfertas []models.PaqueteOferta
	//Itera a través de los resultados y agrega a la slice
	for rows.Next() {
		var paqueteoferta models.PaqueteOferta
		err := rows.Scan(&paqueteoferta.OrigenDestino, &paqueteoferta.Detalle, &paqueteoferta.FechaInicio, &paqueteoferta.FechaFin, &paqueteoferta.Precio)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al escanear resultados", http.StatusInternalServerError)
			return
		}

		fechaInicio, err := time.Parse("2006-01-02T15:04:05Z", paqueteoferta.FechaInicio)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al parsear la fecha de inicio", http.StatusInternalServerError)
			return
		}

		fechaFin, err := time.Parse("2006-01-02T15:04:05Z", paqueteoferta.FechaFin)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al parsear la fecha de fin", http.StatusInternalServerError)
			return
		}

		// Formatea las fechas como solo la parte de la fecha (sin la hora)
		paqueteoferta.FechaInicio = fechaInicio.Format("2006-01-02")
		paqueteoferta.FechaFin = fechaFin.Format("2006-01-02")

		paquetesOfertas = append(paquetesOfertas, paqueteoferta)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paquetesOfertas); err != nil {
		log.Fatal(err)
		http.Error(w, "Error al convertir a JSON", http.StatusInternalServerError)
	}
}
