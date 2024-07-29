package types

import "time"

type tipoOrden int

const (
	OrdenPorNombre tipoOrden = iota
	OrdenPorValor
	OrdenPorFecha
)

type Orden struct {
	TipoOrden  tipoOrden
	Ascendente bool
}

type Filter struct {
	Monedas       []string
	FechaInicial  time.Time
	FechaFinal    time.Time
	TamPaginas    int
	CantPaginas   int
	Orden         Orden
	PaginaInicial int
	Usuario       int
}

type Resumen map[string]any

func (o Orden) TipoToString() string {
	switch o.TipoOrden {
	case OrdenPorNombre:
		return "nombre"
	case OrdenPorFecha:
		return "fecha"
	case OrdenPorValor:
		return "valor"
	default:
		return ""
	}
}

func (o Orden) DireccionToString() string {
	if o.Ascendente {
		return "ASC"
	}
	return "DESC"
}

func (o Orden) ToString() string {
	return o.TipoToString() + " " + o.DireccionToString()
}
