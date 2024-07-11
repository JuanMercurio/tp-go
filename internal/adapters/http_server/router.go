package http_server

import (
	_ "github.com/juanmercurio/tp-go/docs"

	"github.com/gin-gonic/gin"
	"github.com/juanmercurio/tp-go/internal/adapters/http_server/handlers"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	*gin.Engine
}

func Config(h handlers.MonedaHandler) Router {
	router := gin.Default()
	router.GET("/monedas", h.BuscarTodos)
	router.GET("/cotizacion", h.CotizacionMoneda)
	router.GET("/cotizaciones", h.Cotizaciones)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return Router{router}
}

func (r Router) Start() {
	r.Run()
}
