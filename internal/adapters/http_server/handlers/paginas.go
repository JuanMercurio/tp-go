package handlers

import (
	"math"

	"github.com/juanmercurio/tp-go/internal/ports"
)

type Pagina []any

func crearPaginas(tamPaginas int, elementos []any) []Pagina {
	var paginas []Pagina
	restantes := len(elementos)
	cantidadPaginasFinal := paginasNecesarias(tamPaginas, len(elementos))

	for i := 0; i < cantidadPaginasFinal; i++ {
		offset := tamPaginas * i
		if restantes < tamPaginas {
			tamPaginas = restantes
		}
		p := obtenerPaginas(elementos, offset, tamPaginas)
		paginas = append(paginas, p)
		restantes -= tamPaginas
	}

	return paginas
}
func paginasNecesarias(tamPaginas int, filas int) int {
	if filas < tamPaginas {
		return 1
	}
	return int(math.Ceil(float64(filas) / float64(tamPaginas)))
}

func obtenerPaginas(cotizaciones []any, offset int, cant int) Pagina {
	return cotizaciones[offset : offset+cant]
}

func cotizacionesToAny(cotizaciones []ports.CotizacionOutputDTO) []any {
	s := make([]any, len(cotizaciones))
	for i, cotizacion := range cotizaciones {
		s[i] = cotizacion
	}
	return s
}
