package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/juanmercurio/tp-go/internal/adapters/config"
	"github.com/juanmercurio/tp-go/internal/adapters/cotizador"
	"github.com/juanmercurio/tp-go/internal/adapters/cotizador/coinbase"
	"github.com/juanmercurio/tp-go/internal/adapters/cotizador/paprika"
	"github.com/juanmercurio/tp-go/internal/adapters/http_server"
	"github.com/juanmercurio/tp-go/internal/adapters/http_server/handlers"
	"github.com/juanmercurio/tp-go/internal/adapters/mysql"
	"github.com/juanmercurio/tp-go/internal/adapters/mysql/repos"
	"github.com/juanmercurio/tp-go/internal/core/servicios"
)

// @title			Criptomoneda API
// @description	API en la cual podemos consultar cotizaciones de diferentes monedas
func main() {

	config, err := config.Crear()
	if err != nil {
		log.Fatal("Error de config: ", err)
	}

	clienteSQL, err := mysql.CrearCliente(&config)
	if err != nil {
		log.Fatal("Error al inicial el cliente mySQL: ", err)
	}

	repoMonedas := repos.CrearRepositorioMoneda(clienteSQL)

	paprikaAPI := paprika.Crear(&config.Apis.Paprika)
	coinbaseAPI := coinbase.Crear(&config.Apis.CoinBase)
	cotizador := cotizador.Crear(paprikaAPI, coinbaseAPI)

	servicio := servicios.CrearServicioMoneda(repoMonedas, cotizador)

	handlerMoneda := handlers.CrearHandlerMoneda(servicio)

	server := http_server.Config(handlerMoneda)

	server.Start()
}
