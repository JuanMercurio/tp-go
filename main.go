package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// type criptomoneda struct {
//    nombre string
// }

type criptomoneda struct {
	id     int
	nombre string
}
type cotizacion struct {
	moneda criptomoneda
	valor  float64
	time   time.Time
}

var monedas = []criptomoneda{
	{1, "bitcoin"}, {2, "etherium"},
}

func getMonedas(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, monedas)
}

func main() {

	godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	// dbName := os.Getenv("DB_NAME")

	// dns := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", dbUser, dbPass, dbHost, dbPort, dbName)
	dns := fmt.Sprintf("%v:%v@tcp(%v:%v)/", dbUser, dbPass, dbHost, dbPort)

	db, err := sql.Open("mysql", dns)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Conexion correcta a la base de datos")

	// _, err1 := db.Query("CREATE TABLE if not exists numero(id int(11) primary key);")
	// if err1 != nil {
	// 	log.Fatal(err1)
	// }

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

	// router := gin.Default()
	// router.GET("/monedas", getMonedas)
	// router.Run()

}
