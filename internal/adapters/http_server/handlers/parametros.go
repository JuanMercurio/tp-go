package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juanmercurio/tp-go/internal/ports"
)

func validarParametros(c *gin.Context) (ports.Filter, error) {
	fechaInicial, fechaFinal, err := validarFechas(c.Query("fecha_inicial"), c.Query("fecha_final"))
	if err != nil {
		return ports.Filter{}, fmt.Errorf("error en la validacion de fechas: %w", err)
	}

	tamPaginas, err := strconv.Atoi(c.DefaultQuery("tam_paginas", "50"))
	if err != nil {
		return ports.Filter{}, fmt.Errorf("el tamaño de las paginas debe ser un numero entero: %w", err)
	}

	cantPaginas, err := strconv.Atoi(c.DefaultQuery("cant_paginas", "10"))
	if err != nil {
		return ports.Filter{}, fmt.Errorf("a cantidad de paginas debe ser un numero entero: %w", err)
	}

	orden := validarOrden(c.Query("orden"), c.Query("orden_direccion"))

	///TODO falta validar el valor de retorno de error
	paginaInicial, _ := strconv.Atoi(c.DefaultQuery("pagina_inicial", "1"))
	paginaInicial--

	usuarioId, err := strconv.Atoi(c.DefaultQuery("usuario", "0"))
	if err != nil {
		return ports.Filter{}, fmt.Errorf("el id de usuario debe ser un numero")
	}

	var monedasSlice []string
	if monedas := c.Query("monedas"); monedas != "" {
		monedasSlice = strings.Split(monedas, " ")
	}

	parametros := ports.Filter{
		FechaInicial:  fechaInicial,
		FechaFinal:    fechaFinal,
		TamPaginas:    min(50, tamPaginas), // este valor y el de abajo esta harcodeados TODO
		CantPaginas:   min(10, cantPaginas),
		Orden:         orden,
		PaginaInicial: paginaInicial,
		Monedas:       monedasSlice,
		Usuario:       usuarioId,
	}

	// TODO cambiar los valores harcodeados de maximos
	return parametros, nil
}

func validarOrden(orden string, ordenDireccion string) ports.Orden {
	filtro := ports.Orden{}
	switch orden {
	case "nombre":
		filtro.TipoOrden = ports.OrdenPorNombre
	case "valor":
		filtro.TipoOrden = ports.OrdenPorValor
	case "fecha":
		filtro.TipoOrden = ports.OrdenPorFecha
	default:
		filtro.TipoOrden = ports.OrdenPorFecha
	}

	filtro.Ascendente = ordenDireccion == "asc" || ordenDireccion == "ascendente"

	return filtro
}

func validarFechas(fechaInicial string, fechaFinal string) (time.Time, time.Time, error) {

	if fechaFinal == "" && fechaInicial == "" {
		return time.Time{}, time.Time{}, nil
	}

	var fechaInicialValida time.Time
	var fechaFinalValida time.Time
	var errs []error
	var err error

	if fechaFinal != "" {
		fechaFinalValida, err = validarFecha(fechaFinal)
		if err != nil {
			errs = append(errs, fmt.Errorf("error de validacion de fecha_inicial:  %w", err))
		}
	}

	if fechaFinal != "" {
		fechaFinalValida, err = validarFecha(fechaFinal)
		if err != nil {
			errs = append(errs, fmt.Errorf("error de validacion de fecha_final:  %w", err))
		}
	}

	if !rangoValido(fechaInicialValida, fechaFinalValida) {
		errs = append(errs, fmt.Errorf("el rango de las fechas es invalido"))
	}

	if len(errs) > 0 {
		return time.Time{}, time.Time{}, errors.Join(errs...)
	}

	return fechaInicialValida, fechaFinalValida, nil
}

func rangoValido(inicial, final time.Time) bool {
	if inicial.IsZero() || final.IsZero() {
		return true
	}

	if inicial.After(final) || inicial == final {
		return false
	}

	return true
}

func validarFecha(fechaObtenida string) (time.Time, error) {

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
