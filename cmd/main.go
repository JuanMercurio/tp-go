package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/juanmercurio/tp-go/internal/adapters/config"
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
	repoUsuarios := repos.CrearRepositorioUsuario(clienteSQL)

	paprikaAPI := paprika.Crear(&config.Apis.Paprika)
	coinbaseAPI := coinbase.Crear(&config.Apis.CoinBase)

	servicioMoneda := servicios.CrearServicioMoneda(repoMonedas, paprikaAPI, coinbaseAPI)
	servicioUsuario := servicios.CrearServicioUsuario(repoUsuarios, repoMonedas)

	handlerMoneda := handlers.CrearHandlerMoneda(servicioMoneda)
	handlerUsuario := handlers.CrearHandlerUsuario(servicioUsuario)

	server := http_server.Config(handlerMoneda, handlerUsuario)

	server.Start()

	//todo end
}
