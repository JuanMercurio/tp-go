package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func CrearCliente() *sql.DB {
	//TODO encontrar un mejor lugar para inicializar las variables de entorne
	//TODO analizar si log Fatal o retornal el error

	godotenv.Load()
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dns := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dns)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	log.Println("Conexion correcta a la base de datos")

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var s []uint8
	for rows.Next() {
		if err := rows.Scan(&s); err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(s))
	}

	return db
}
