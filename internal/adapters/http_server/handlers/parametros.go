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

func validarParametros(c *gin.Context) (ports.ParamCotizaciones, error) {
	fechaInicial, fechaFinal, err := validarFechas(c.Query("fecha_inicial"), c.Query("fecha_final"))
	if err != nil {
		return ports.ParamCotizaciones{}, fmt.Errorf("error en la validacion de fechas: %w", err)
	}

	tamPaginas, err := strconv.Atoi(c.DefaultQuery("tam_paginas", "50"))
	if err != nil {
		return ports.ParamCotizaciones{}, fmt.Errorf("el tamaño de las paginas debe ser un numero entero: %w", err)
	}

	cantPaginas, err := strconv.Atoi(c.DefaultQuery("cant_paginas", "1"))
	if err != nil {
		return ports.ParamCotizaciones{}, fmt.Errorf("a cantidad de paginas debe ser un numero entero: %w", err)
	}

	orden := validarOrden(c.Query("orden"), c.Query("orden_direccion"))

	///TODO falta validar el valor de retorno de error
	paginaInicial, _ := strconv.Atoi(c.DefaultQuery("pagina_inicial", "1"))
	paginaInicial--

	monedasSlice := strings.Split(c.Query("monedas"), " ")

	parametros := ports.ParamCotizaciones{
		FechaInicial:  fechaInicial,
		FechaFinal:    fechaFinal,
		TamPaginas:    min(50, tamPaginas), // este valor y el de abajo esta harcodeados TODO
		CantPaginas:   min(10, cantPaginas),
		Orden:         orden,
		PaginaInicial: paginaInicial,
		Monedas:       monedasSlice,
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
