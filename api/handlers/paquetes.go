package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"backend/models"
	"backend/utils"
)

func ObtenerPaquetes(w http.ResponseWriter, r *http.Request) {
	// Captura los parámetros de la URL
	origenIDStr := r.URL.Query().Get("origen")
	destinoIDStr := r.URL.Query().Get("destino")
	fechaInicio := r.URL.Query().Get("fechaInicio")
	fechaFin := r.URL.Query().Get("fechaFin")

	// Convierte los ID de origen y destino a enteros
	origenID, err := strconv.Atoi(origenIDStr)
	if err != nil {
		http.Error(w, "El parámetro 'origen' debe ser un número válido", http.StatusBadRequest)
		return
	}

	destinoID, err := strconv.Atoi(destinoIDStr)
	if err != nil {
		http.Error(w, "El parámetro 'destino' debe ser un número válido", http.StatusBadRequest)
		return
	}

	// Abre la conexión con la base de datos
	db, err := utils.OpenDB()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Realiza la consulta SQL con los parámetros
	rows, err := db.Query(`
	SELECT 
            fp.id AS id_paquete,
            p.nombre AS nombre_paquete,
            p.descripcion AS desc_paquete,
            p.detalles AS detalle_paquete,
            fp.fechaInit,
            fp.fechaFin,
            (fp.fechaFin - fp.fechaInit) as cant_dias,
            p.precio_viaje AS pr_viaje,
            hh.precio_noche AS pr_noche,
            (p.precio_viaje + (hh.precio_noche * (fp.fechaFin - fp.fechaInit))) as pr_total,
            co.nombre AS ciudad_origen,
            cd.nombre AS ciudad_destino,
            h.nombre AS nombre_hotel,
            hh.descripcion AS desc_hh,
            hh.servicios AS servicios_hh
        FROM 
            paquete AS p
        JOIN 
            aeropuerto AS ao ON p.id_origen = ao.id
        JOIN 
            aeropuerto AS ad ON p.id_destino = ad.id
        JOIN 
            Ciudad AS co ON ao.ciudad_id = co.id
        JOIN 
            ciudad AS cd ON ad.ciudad_id = cd.id
        JOIN 
            habitacionhotel AS hh ON p.id_hh = hh.id
        JOIN 
            hotel AS h ON hh.hotel_id = h.id
        JOIN 
            fechaPaquete AS fp ON p.id = fp.id_paquete
        WHERE 
            ao.id = $1
            AND ad.id = $2
            AND fp.fechaInit = $3
            AND fp.fechaFin = $4;

    `, origenID, destinoID, fechaInicio, fechaFin)
	if err != nil {
		http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Procesa los resultados y crea una estructura para la respuesta JSON
	var paquetes []models.Paquete // Utiliza models.Paquete
	for rows.Next() {
		var paquete models.Paquete // Utiliza models.Paquete
		if err := rows.Scan(
			&paquete.ID,
			&paquete.Nombre,
			&paquete.Descripcion,
			&paquete.Detalles,
			&paquete.FechaInit,
			&paquete.FechaFin,
			&paquete.CantDias,
			&paquete.PrecioViaje,
			&paquete.PrecioNoche,
			&paquete.PrecioTotal,
			&paquete.CiudadOrigen,
			&paquete.CiudadDestino,
			&paquete.NombreHotel,
			&paquete.DescripcionHH,
			&paquete.ServiciosHH,
		); err != nil {
			http.Error(w, "Error al escanear resultados", http.StatusInternalServerError)
			return
		}

		fechaInicio, err := time.Parse("2006-01-02T15:04:05Z", paquete.FechaInit)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al parsear la fecha de inicio", http.StatusInternalServerError)
			return
		}

		fechaFin, err := time.Parse("2006-01-02T15:04:05Z", paquete.FechaFin)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al parsear la fecha de fin", http.StatusInternalServerError)
			return
		}

		// Formatea las fechas como solo la parte de la fecha (sin la hora)
		paquete.FechaInit = fechaInicio.Format("2006-01-02")
		paquete.FechaFin = fechaFin.Format("2006-01-02")

		paquetes = append(paquetes, paquete)
	}

	// Convierte los resultados a JSON y envía la respuesta
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paquetes); err != nil {
		http.Error(w, "Error al convertir a JSON", http.StatusInternalServerError)
	}
}

func ObtenerPaquetesMes(w http.ResponseWriter, r *http.Request) {
	// Captura los parámetros de la URL
	origenIDStr := r.URL.Query().Get("origen")
	destinoIDStr := r.URL.Query().Get("destino")
	mesStr := r.URL.Query().Get("mes")

	// Convierte los ID de origen, destino y mes a enteros
	origenID, err := strconv.Atoi(origenIDStr)
	if err != nil {
		http.Error(w, "El parámetro 'origen' debe ser un número válido", http.StatusBadRequest)
		return
	}

	destinoID, err := strconv.Atoi(destinoIDStr)
	if err != nil {
		http.Error(w, "El parámetro 'destino' debe ser un número válido", http.StatusBadRequest)
		return
	}

	mes, err := strconv.Atoi(mesStr)
	if err != nil {
		http.Error(w, "El parámetro 'mes' debe ser un número válido", http.StatusBadRequest)
		return
	}

	// Abre la conexión con la base de datos
	db, err := utils.OpenDB()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Realiza la consulta SQL con los parámetros
	rows, err := db.Query(`
        SELECT 
            fp.id AS id_paquete,
            p.nombre AS nombre_paquete,
            p.descripcion AS desc_paquete,
            p.detalles AS detalle_paquete,
            fp.fechaInit,
            fp.fechaFin,
            (fp.fechaFin - fp.fechaInit) as cant_dias,
            p.precio_viaje AS pr_viaje,
            hh.precio_noche AS pr_noche,
            (p.precio_viaje + (hh.precio_noche * (fp.fechaFin - fp.fechaInit))) as pr_total,
            co.nombre AS ciudad_origen,
            cd.nombre AS ciudad_destino,
            h.nombre AS nombre_hotel,
            hh.descripcion AS desc_hh,
            hh.servicios AS servicios_hh
        FROM 
            paquete AS p
        JOIN 
            aeropuerto AS ao ON p.id_origen = ao.id
        JOIN 
            aeropuerto AS ad ON p.id_destino = ad.id
        JOIN 
            Ciudad AS co ON ao.ciudad_id = co.id
        JOIN 
            ciudad AS cd ON ad.ciudad_id = cd.id
        JOIN 
            habitacionhotel AS hh ON p.id_hh = hh.id
        JOIN 
            hotel AS h ON hh.hotel_id = h.id
        JOIN 
            fechaPaquete AS fp ON p.id = fp.id_paquete
        WHERE 
            ao.id = $1
            AND ad.id = $2
            AND EXTRACT(MONTH FROM fp.fechaInit) = $3  -- Filtra por el mes especificado
            AND EXTRACT(MONTH FROM fp.fechaFin) = $3   -- Filtra por el mes especificado
    `, origenID, destinoID, mes)
	if err != nil {
		http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Procesa los resultados y crea una estructura para la respuesta JSON
	var paquetes []models.Paquete
	for rows.Next() {
		var paquete models.Paquete // Utiliza models.Paquete
		if err := rows.Scan(
			&paquete.ID,
			&paquete.Nombre,
			&paquete.Descripcion,
			&paquete.Detalles,
			&paquete.FechaInit,
			&paquete.FechaFin,
			&paquete.CantDias,
			&paquete.PrecioViaje,
			&paquete.PrecioNoche,
			&paquete.PrecioTotal,
			&paquete.CiudadOrigen,
			&paquete.CiudadDestino,
			&paquete.NombreHotel,
			&paquete.DescripcionHH,
			&paquete.ServiciosHH,
		); err != nil {
			http.Error(w, "Error al escanear resultados", http.StatusInternalServerError)
			return
		}

		fechaInicio, err := time.Parse("2006-01-02T15:04:05Z", paquete.FechaInit)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al parsear la fecha de inicio", http.StatusInternalServerError)
			return
		}

		fechaFin, err := time.Parse("2006-01-02T15:04:05Z", paquete.FechaFin)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al parsear la fecha de fin", http.StatusInternalServerError)
			return
		}

		// Formatea las fechas como solo la parte de la fecha (sin la hora)
		paquete.FechaInit = fechaInicio.Format("2006-01-02")
		paquete.FechaFin = fechaFin.Format("2006-01-02")

		paquetes = append(paquetes, paquete)
	}

	// Convierte los resultados a JSON y envía la respuesta
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paquetes); err != nil {
		http.Error(w, "Error al convertir a JSON", http.StatusInternalServerError)
	}
}
