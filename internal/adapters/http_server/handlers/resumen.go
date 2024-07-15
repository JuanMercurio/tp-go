package handlers

import "github.com/juanmercurio/tp-go/internal/ports"

type Filtros struct {
	Monedas        string `form:"monedas"`
	FechaInicial   string `form:"fecha_inicial"`
	FechaFinal     string `form:"fecha_final"`
	TamPaginas     int    `form:"tam_paginas"`
	PaginaInicial  int    `form:"pagina_inicial"`
	CantPaginas    int    `form:"cant_paginas"`
	Orden          string `form:"orden"`
	OrdenDireccion string `form:"orden_direccion"`
}

func crearResumen(params ports.ParamCotizaciones) map[string]any {
	filtros := make(map[string]any)
	filtros["monedas"] = params.Monedas
	filtros["fecha_inicial"] = params.FechaInicial
	filtros["fecha_final"] = params.FechaFinal
	filtros["tam_paginas"] = params.TamPaginas
	filtros["cant_paginas"] = params.CantPaginas
	filtros["orden"] = params.Orden.TipoToString()
	filtros["orden_direccion"] = params.Orden.DireccionToString()
	filtros["pagina_inicial"] = params.PaginaInicial + 1
	return filtros
}
