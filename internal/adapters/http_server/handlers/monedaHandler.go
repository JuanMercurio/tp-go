package handlers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/core/servicios"
	"github.com/juanmercurio/tp-go/internal/ports"
)

type Pagina []any

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
	}

	c.JSON(http.StatusOK, todos)
}

// @Summary		Retorna las cotizaciones paginadas segun los filtros
// @Description	Podemos filtar segun nombre de moneda, fecha, y cantidad de registros.
// @Tags			Moneda
// @Accept			json
// @Produce		json
// @Param			monedas			query		string	true	"id de las monedas que queremos separados por espacios"
// @Param			fecha_inicial	query		string	false	"Fecha desde la cual quiero obtener cotizaciones"
// @Param			fecha_final		query		string	false	"Fecha hasta la cual quiero obtener cotizaciones"
// @Param			tam_paginas		query		int		false	"El tamaño de las paginas, como maximo es 100, el default es 50"
// @Param			cant_paginas	query		int		false	"La cantidad de paginas, como maximo es 10, el default es 1"
// @Success		200				{object}	[]Pagina
// @Failure		400				{object}	string
// @Router			/cotizaciones [get]
func (mh MonedaHandler) Cotizaciones(c *gin.Context) {
	// TODO cambiar a que retorna siempre 10 paginas y recibe un offset de paginas

	// ojo con la cantidad de returns (struct?)
	parametros, err := validarParametros(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	monedasSlice := strings.Split(c.Query("monedas"), " ")
	parametros.Monedas = monedasSlice

	cotizaciones, err := mh.srv.Cotizaciones(parametros)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	paginas := crearPaginas(parametros.TamPaginas, cotizaciones)
	c.JSON(http.StatusOK, paginas)
}

// @Summary		Retorna las cotizaciones de una moneda
// @Description	Segun una moneda, una fecha de inicio y fin, y una cantidad, nos retorna las cotizaciones
// @Tags			Moneda
// @Accept			json
// @Produce		json
// @Param			moneda_id		query		string	true	"El ID de la moneda que queremos"
// @Param			fecha_inicial	query		string	true	"Fecha desde la cual quiero obtener cotizaciones"
// @Param			fecha_final		query		string	true	"Fecha hasta la cual quiero obtener cotizaciones"
// @Param			offset			query		string	false	"Desde que elemento quiero obtener resutados"
// @Param			cant			query		string	false	"Cantidad de resultados que quiero"
// @Success		200				{object}	[]domain.Cotizacion
// @Success		401				{object}	string
// @Router			/cotizacion [get]
func (mh MonedaHandler) CotizacionMoneda(c *gin.Context) {

	// TODO falta manejar errores aca, fijarse lo de multiples errores
	fechaInicial, _ := time.Parse("2006-01-02 15:04:05", c.Query("fecha_inicial"))
	fechaFinal, _ := time.Parse("2006-01-02 15:04:05", c.Query("fecha_final"))
	monedaID, _ := strconv.Atoi(c.Query("moneda_id"))
	// TODO es el default que nos va a traer si el usuario no indica cantidad
	// podria ser una variable de entorno en vez de un valor hardcodeado
	cant := 10
	if cantString := c.Query("cant"); cantString != "" {
		cant, _ = strconv.Atoi(cantString)
	}
	var offset int
	if offsetString := c.Query("offset"); offsetString != "" {
		offset, _ = strconv.Atoi(offsetString)
	}

	cotizaciones, err := mh.srv.BuscarCotizacionMoneda(monedaID, fechaInicial, fechaFinal, cant, offset)
	if err != nil {
		c.JSON(http.StatusConflict, err)
	}

	c.JSON(http.StatusOK, cotizaciones)
}

func crearPaginas(tamPaginas int, cotizaciones []domain.Cotizacion) []Pagina {
	cotizacionesAny := deCotizacionAAny(cotizaciones)
	var paginas []Pagina
	restantes := len(cotizaciones)
	cantidadPaginasFinal := paginasNecesarias(tamPaginas, len(cotizaciones))

	// for i := 0; i < cantPaginas && i < len(cotizaciones); i++ {
	for i := 0; i < cantidadPaginasFinal; i++ {
		offset := tamPaginas * i
		if restantes < tamPaginas {
			tamPaginas = restantes
		}
		// if offset+tamPaginas > len(cotizaciones) {
		// 	break
		// }
		p := obtenerPagina(cotizacionesAny, offset, tamPaginas)
		paginas = append(paginas, p)
		restantes -= tamPaginas
	}

	return paginas
}
func paginasNecesarias(tamPaginas int, cantCotizaciones int) int {
	if cantCotizaciones < tamPaginas {
		return 1
	}
	return int(math.Ceil(float64(cantCotizaciones) / float64(tamPaginas)))
}

func obtenerPagina(cotizaciones []any, offset int, cant int) Pagina {
	return cotizaciones[offset : offset+cant]
}

func deCotizacionAAny(cotizaciones []domain.Cotizacion) []any {
	s := make([]any, len(cotizaciones))
	for i, cotizacion := range cotizaciones {
		s[i] = cotizacion
	}
	return s

}

func validarParametros(c *gin.Context) (ports.Parametros, error) {
	fechaInicial, fechaFinal, err := validarFechas(c.Query("fecha_inicial"), c.Query("fecha_final"))
	if err != nil {
		return ports.Parametros{}, fmt.Errorf("error en la validacion de fechas: %w", err)
	}

	tamPaginas, err := strconv.Atoi(c.DefaultQuery("tam_paginas", "5"))
	if err != nil {
		return ports.Parametros{}, fmt.Errorf("error en el formato del tamaño de las paginas: %w", err)
	}

	cantPaginas, err := strconv.Atoi(c.DefaultQuery("cant_paginas", "1"))
	if err != nil {
		return ports.Parametros{}, fmt.Errorf("error en el formato de la cantidad de paginas: %w", err)
	}

	filtro := validarFiltro(c.Query("orden"), c.Query("ascendente"))

	parametros := ports.Parametros{
		FechaInicial: fechaInicial,
		FechaFinal:   fechaFinal,
		TamPaginas:   min(50, tamPaginas), // este valor y el de abajo esta harcodeados TODO
		CantPaginas:  min(10, cantPaginas),
		Filtro:       filtro,
	}

	// TODO cambiar los valores harcodeados de maximos
	return parametros, nil
}

func validarFiltro(orden string, ascendente string) ports.Filtro {
	filtro := ports.Filtro{}
	switch orden {
	case "nombre":
		filtro.TipoFiltro = ports.FiltroPorNombre
	case "valor":
		filtro.TipoFiltro = ports.FiltroPorValor
	case "fecha":
		filtro.TipoFiltro = ports.FiltroPorFecha
	default:
		filtro.TipoFiltro = ports.FiltroPorFecha
	}

	//por default es descentente
	if ascendente == "1" {
		filtro.Ascendente = true
	} else {
		filtro.Ascendente = false
	}

	return filtro
}

func validarFechas(fechaInicial string, fechaFinal string) (time.Time, time.Time, error) {

	var errs []error
	fechaInicialDefault := time.Now().AddDate(0, 0, -14)
	fechaFinalDefault := time.Now()

	fechaInicialValidada, err := validarFecha(fechaInicial, fechaInicialDefault)
	if err != nil {
		errs = append(errs, fmt.Errorf("error de validacion de fecha_inicial:  %w", err))
	}

	fechaFinalValidada, err := validarFecha(fechaFinal, fechaFinalDefault)
	if err != nil {
		errs = append(errs, fmt.Errorf("error de validacion de fecha_final: %w", err))
	}

	if fechaInicialValidada.After(fechaFinalValidada) || fechaFinalValidada == fechaInicialValidada {
		return time.Time{}, time.Time{}, fmt.Errorf("el rango de las fechas es invalido")
	}

	if len(errs) > 0 {
		return time.Time{}, time.Time{}, errors.Join(errs...)
	}

	return fechaInicialValidada, fechaFinalValidada, nil
}

func validarFecha(fechaObtenida string, fechaDefault time.Time) (time.Time, error) {
	if fechaObtenida == "" {
		return fechaDefault, nil
	}

	fecha, err := time.Parse("2006-01-02 15:04:05", fechaObtenida)
	if err != nil {
		return time.Time{}, fmt.Errorf("error en el formato de la fecha: %w", err)
	}
	return fecha, nil
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
