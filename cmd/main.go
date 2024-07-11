package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/juanmercurio/tp-go/internal/adapters/http_server"
	"github.com/juanmercurio/tp-go/internal/adapters/http_server/handlers"
	"github.com/juanmercurio/tp-go/internal/adapters/mysql"
	"github.com/juanmercurio/tp-go/internal/adapters/mysql/repos"
	"github.com/juanmercurio/tp-go/internal/core/servicios"
)

// @title			Criptomoneda API
// @description	API en la cual podemos consultar cotizaciones de diferentes monedas
func main() {
	clienteSQL := mysql.CrearCliente() //TODO esto se podria meter dentro del repositorio en si
	repoMonedas := repos.CrearRepositorioMoneda(clienteSQL)
	servicio := servicios.CrearServicioMoneda(repoMonedas)
	handlerMoneda := handlers.CrearHandlerMoneda(servicio)
	server := http_server.Config(handlerMoneda)
	server.Start()
}
