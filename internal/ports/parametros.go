package ports

import "time"

type tipoFiltro int

const (
	FiltroPorNombre tipoFiltro = iota
	FiltroPorValor
	FiltroPorFecha
)

type Filtro struct {
	TipoFiltro tipoFiltro
	Ascendente bool
}

type Parametros struct {
	Monedas      []string
	FechaInicial time.Time
	FechaFinal   time.Time
	TamPaginas   int
	CantPaginas  int
	Filtro       Filtro
}
