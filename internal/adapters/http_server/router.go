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

func Config(mh *handlers.MonedaHandler, uh *handlers.UsuarioHandler) Router {

	router := gin.Default()

	router.GET("/usuarios", uh.BuscarUsuarios)
	router.POST("/usuarios", uh.AltaUsuario)
	router.DELETE("/usuarios", uh.BajaUsuario)
	router.PATCH("/usuarios/:id", uh.ActualizarUsuario)

	router.POST("/cotizacion", mh.AltaCotizacionManual)
	router.DELETE("/cotizacion/:id", mh.BajaCotizacion)
	router.PATCH("/cotizacion/:id", mh.ActualizarCotizacion)

	router.GET("/monedas", mh.BuscarMonedas)
	router.GET("/cotizaciones", mh.Cotizaciones)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Use(handlers.Autenticar)
	router.POST("/monedas", mh.AltaMoneda)
	router.POST("/cotizaciones", mh.AltaCotizaciones)

	return Router{router}
}

func (r Router) Start() {
	r.Run()
}
