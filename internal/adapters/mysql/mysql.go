package mysql

import (
	"database/sql"
	"fmt"

	"github.com/juanmercurio/tp-go/internal/adapters/config"
)

func CrearCliente(config *config.Config) (*sql.DB, error) {

	dbUser := config.ENV["DB_USER"]
	dbPass := config.ENV["DB_PASS"]
	dbHost := config.ENV["DB_HOST"]
	dbPort := config.ENV["DB_PORT"]
	dbName := config.ENV["DB_NAME"]

	dns := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dns)
	if err != nil {
		return db, err
	}

	if err := db.Ping(); err != nil {
		return db, err
	}

	return db, nil
}
