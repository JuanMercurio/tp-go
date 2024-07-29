package handlers

import "github.com/juanmercurio/tp-go/internal/ports/types"

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

func (h CotizacionHandler) crearResumen(params types.Filter) (types.Resumen, error) {
	resumen, err := h.sc.Resumen(params)
	if err != nil {
		return types.Resumen{}, err
	}

	return resumen, nil
}
