package handlers

import (
	"github.com/juanmercurio/tp-go/internal/ports"
)

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

func (h MonedaHandler) crearResumen(params ports.Filter) (ports.Resumen, error) {
	resumen, err := h.srv.Resumen(params)
	if err != nil {
		return ports.Resumen{}, err
	}

	return resumen, nil
}
