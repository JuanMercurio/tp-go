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

func Config(h *handlers.MonedaHandler) Router {

	router := gin.Default()

	router.POST("/monedas", handlers.Autenticar, h.AltaMoneda)
	router.POST("/cotizaciones", handlers.Autenticar, h.AltaCotizaciones)

	router.POST("/usuario", h.AltaUsuario)
	router.DELETE("/usuario", h.BajaUsuario)

	router.GET("/monedas", h.BuscarTodos)
	router.GET("/cotizaciones", h.Cotizaciones)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return Router{router}
}

func (r Router) Start() {
	r.Run()
}
