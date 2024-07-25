package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juanmercurio/tp-go/internal/ports"
)

type MonedaHandler struct {
	srv ports.ServicioMonedas
}

func CrearHandlerMoneda(srv ports.ServicioMonedas) *MonedaHandler {
	return &MonedaHandler{
		srv: srv,
	}
}

// @Summary		Busca todas las monedas
// @Description	Obtiene una lista de todos las monedas disponibles.
// @Tags			Moneda
// @Accept			json
// @Produce		json
// @Success		200	{object}	[]ports.MonedaOutputDTO
// @Router			/monedas [get]
func (mh MonedaHandler) BuscarMonedas(c *gin.Context) {
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
// @Param			simbolo			query		string	true	"Simbolo de la moneda"
// @Param			nombre			query		string	false	"Nombre de la moneda nueva"
// @Success		200				{object}	int
// @Failure		400				{object}	error
// @Router			/monedas [post]
func (mh MonedaHandler) AltaMoneda(c *gin.Context) {
	nombre := c.Query("nombre")
	simbolo := c.Query("simbolo")

	id, err := mh.srv.AltaMoneda(nombre, simbolo)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	mh.srv.CotizarNuevaMoneda(id, simbolo)

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary	Retorna las cotizaciones paginadas segun los filtros
// @Tags		Cotizacion
// @Accept		json
// @Produce	json
// @Param		monedas			query		string	false	"simbolos de las monedas que queremos separados por espacios"
// @Param		fecha_inicial	query		string	false	"Fecha desde la cual quiero obtener cotizaciones (YYYY-MM-DD HH:MM:SS)"	Format(date-time)
// @Param		fecha_final		query		string	false	"Fecha hasta la cual quiero obtener cotizaciones (YYYY-MM-DD HH:MM:SS)"	Format(date-time)
// @Param		tam_paginas		query		string	false	"El tamaño de las paginas, como maximo es 100, el default es 50"
// @Param		pagina_inicial	query		int		false	"Pagina a partir de la cual sera retornado el query"
// @Param		cant_paginas	query		int		false	"La cantidad de paginas, como maximo es 10, el default es 10"
// @Param		orden			query		string	false	"Ordenar segun alguno de estos valores: fecha(default), valor, nombre"	Enum(fecha, valor, nombre)
// @Param		orden_direccion	query		string	false	"Indica si es ascendente o descendente, el default es desdencente"		Enum(ascendente, descendente)
// @Param		usuario			query		int		false	"Usuario elegido"
// @Param		resumen			query		string	false	"Para incluir resumen indicar el valor debe ser si"	Enum(si, no)
// @Success	200				{object}	[]Pagina
// @Failure	400				{object}	string
// @Router		/cotizaciones [get]
func (mh MonedaHandler) Cotizaciones(c *gin.Context) {

	// esto nos devuelve un map de todos los query params
	// filtros := c.Request.URL.Query()

	parametrosBusqueda, err := validarParametros(c)
	quiereResumen := c.Query("resumen") == "si"

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cantidad, cotizaciones, err := mh.srv.Cotizaciones(parametrosBusqueda)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	paginas := crearPaginas(parametrosBusqueda.TamPaginas, cotizacionesToAny(cotizaciones))
	body := gin.H{"cotizaciones": paginas}

	if quiereResumen {
		body["resumen"], err = mh.crearResumen(parametrosBusqueda)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "no se pudo crear el resumen"})
		}
		body["cantidad_total_de_registros"] = cantidad
		body["cantidad_de_paginas_posibles"] = cantidad / parametrosBusqueda.TamPaginas
	}

	c.JSON(http.StatusOK, body)
}

// @Summary	Llama para que se haga la cotizacion de las monedas
// @Tags		Cotizacion
// @Accept		json
// @Produce	json
// @Param		Authorization	header		string	true	"Token de autorización"
// @Param		api				query		string	true	"API que se utilizara para cotizar"
// @Success	200				{object}	[]Pagina
// @Failure	400				{object}	string
// @Router		/cotizaciones [post]
func (mh MonedaHandler) AltaCotizaciones(c *gin.Context) {

	api := c.Query("api")
	if !ApiValida(api) {
		c.JSON(http.StatusConflict, gin.H{
			"error": fmt.Sprintf("no soportamos la api %s, las validas son: Paprika y CoinBase", api),
		})
		return
	}

	responseBody := make(map[string]any)
	err := mh.srv.AltaCotizaciones(api)

	if err != nil {
		errores := strings.Split(err.Error(), "\n")
		responseBody["errores"] = errores
		c.JSON(http.StatusConflict, responseBody)
		return
	}

	c.JSON(http.StatusOK, "Se realizo con exito la cotizacion")
}

func ApiValida(nombre string) bool {
	return nombre == "Paprika" || nombre == "CoinBase"
}

// @Summary	Usuario hace cotizacion de moneda manualmente
// @Tags		Cotizacion
// @Accept		json
// @Produce	json
// @Param		id-usuario	query		string	true	"Usuario que cotizara"
// @Param		simbolo		query		string	true	"Simbolo de la moneda que cotizara"
// @Param		precio		query		string	true	"Valor que cotizara"
// @Param		fecha		query		string	true	"Fecha de la cotizacion"
// @Success	200			{object}	string
// @Failure	400			{object}	string
// @Router		/cotizacion [post]
func (mh MonedaHandler) AltaCotizacionManual(c *gin.Context) {
	fecha, err := stringAFecha(c.Query("fecha"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// idUsuario := c.Query("id-usuario")

	precio, err := strconv.ParseFloat(c.Query("precio"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el precio no es valido"})
		return
	}

	simbolo := c.Query("simbolo")

	if err := mh.srv.SimboloValido(simbolo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("el simbolo no es valido: %w", err).Error()})
		return
	}

	err = mh.srv.CotizarManualmente(simbolo, fecha, precio)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Se realizo con exito la cotizacion")
}

// @Summary	Usuario elimina cotizacion de moneda manualmente
// @Tags		Cotizacion
// @Accept		json
// @Produce	json
// @Param		id				path		string	true	"Usuario que elimina"
// @Param		id-cotizacion	query		string	true	"cotizacion a eliminar"
// @Success	200				{object}	string
// @Failure	400				{object}	string
// @Router		/cotizacion/{id} [delete]
func (mh MonedaHandler) BajaCotizacion(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id-cotizacion"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el id debe ser un entero"})
		return
	}
	err = mh.srv.BajaCotizacionManual(id)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Se elimino con exito la cotizacion")
}

// @Summary	Usuario cambia cotizacion de moneda manualmente
// @Tags		Cotizacion
// @Accept		json
// @Produce	json
// @Param		id			path		string			true	"id de cotizacion"
// @Param		id-usuario	query		string			true	"Usuario que hace los cambios"
// @Param		cambios		body		[]ports.Patch	true	"Usuario que hace los cambios"
// @Success	200			{object}	string
// @Failure	400			{object}	string
// @Router		/cotizacion/{id} [patch]
func (mh MonedaHandler) ActualizarCotizacion(c *gin.Context) {

	var cambios []ports.Patch
	if err := c.BindJSON(&cambios); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("formato incorrecto: %w", err).Error()})
		return
	}

	idCotizacion, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el id de la cotizacion debe ser un entero"})
		return
	}

	idUsuario, err := strconv.Atoi(c.Query("id-usuario"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el id del usuario debe ser un entero"})
		return
	}

	err = mh.srv.PatchCotizacion(idUsuario, idCotizacion, cambios)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Se actualizo con exito la cotizacion")
}
