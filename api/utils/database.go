package utils

import (
	"backend/api/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// OpenDB abre la conexión con la base de datos y la devuelve
func OpenDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
