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

func (h MonedaHandler) crearResumen(params ports.Filter) map[string]any {
	filtros := make(map[string]any)
	// filtros["monedas_disponibles"] = h.srv.MonedasDisponibles(params)
	// filtros["fecha_inicial"] = h.srv.FechaMinima(params)
	// filtros["fecha_final"] = h.srv.FechaMaxima(params)
	// filtros["orden"] = params.Orden.TipoToString()
	// filtros["orden_direccion"] = params.Orden.DireccionToString()
	return filtros
}
