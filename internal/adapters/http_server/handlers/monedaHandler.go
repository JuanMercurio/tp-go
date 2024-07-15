package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanmercurio/tp-go/internal/core/servicios"
)

type MonedaHandler struct {
	srv servicios.MonedaServicio
}

func CrearHandlerMoneda(srv servicios.MonedaServicio) MonedaHandler {
	return MonedaHandler{
		srv: srv,
	}
}

// @Summary		Busca todas las monedas
// @Description	Obtiene una lista de todos las monedas disponibles.
// @Tags			Moneda
// @Accept			json
// @Produce		json
// @Success		200	{object}	[]domain.Criptomoneda
// @Success		401	{object}	string
// @Router			/monedas [get]
func (mh MonedaHandler) BuscarTodos(c *gin.Context) {
	todos, err := mh.srv.BuscarTodos()
	if err != nil {
		c.JSON(http.StatusConflict, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

// @Summary		Da de alta una moneda
// @Description	Si tenemos las credenciales podemos dar de alta una moneda
// @Tags			Moneda
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string	true	"Token de autorización"
// @Param			nombre			query		string	true	"Nombre de la moneda nueva"
// @Success		200				{object}	string
// @Failure		400				{object}	string
// @Router			/monedas [post]
func (mh MonedaHandler) AltaMoneda(c *gin.Context) {
	// TODO autenticacion
	if err := autenticarUsuario(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := mh.srv.AltaMoneda(c.Query("nombre"))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary	Retorna las cotizaciones paginadas segun los filtros
// @Tags		Moneda
// @Accept		json
// @Produce	json
// @Param		monedas			query		string	true	"id de las monedas que queremos separados por espacios"
// @Param		fecha_inicial	query		string	false	"Fecha desde la cual quiero obtener cotizaciones"
// @Param		fecha_final		query		string	false	"Fecha hasta la cual quiero obtener cotizaciones"
// @Param		tam_paginas		query		string	false	"El tamaño de las paginas, como maximo es 100, el default es 50"
// @Param		pagina_inicial	query		int		false	"Pagina a partir de la cual sera retornado el query"
// @Param		cant_paginas	query		int		false	"La cantidad de paginas, como maximo es 10, el default es 1"
// @Param		orden			query		string	false	"El orden en el cual se devuelven las cotizaciones, el default es por fecha"
// @Param		orden_direccion	query		string	false	"Indica si es ascendente o descendente, el default es desdencente"
// @Param		resumen			query		string	false	"Para incluir resumen"
// @Success	200				{object}	[]Pagina
// @Failure	400				{object}	string
// @Router		/cotizaciones [get]
func (mh MonedaHandler) Cotizaciones(c *gin.Context) {

	// esto nos devuelve un map de todos los query params
	// filtros := c.Request.URL.Query()

	parametrosBusqueda, err := validarParametros(c)
	resumen := c.Query("resumen") == "si"

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cotizaciones, err := mh.srv.Cotizaciones(parametrosBusqueda)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	paginas := crearPaginas(parametrosBusqueda.TamPaginas, cotizaciones)
	body := gin.H{"cotizaciones": paginas}

	if resumen {
		body["resumen"] = crearResumen(parametrosBusqueda)
	}

	c.JSON(http.StatusOK, body)
}

// @Summary	Llama para que se haga la cotizacion de las monedas
// @Tags		Moneda
// @Accept		json
// @Produce	json
// @Param		Authorization	header		string	true	"Token de autorización"
// @Success	200				{object}	[]Pagina
// @Failure	400				{object}	string
// @Router		/cotizaciones [post]
func (mh MonedaHandler) AltaCotizaciones(c *gin.Context) {
	// TODO autenticacion
	if err := autenticarUsuario(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	cotizaciones, err := mh.srv.AltaCotizaciones()
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cotizaciones": cotizaciones})
}
