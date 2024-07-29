package handlers

import (
	"math"

	"github.com/juanmercurio/tp-go/internal/ports/types"
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
		p := obtenerElementos(elementos, offset, tamPaginas)
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

func obtenerElementos(elementos []any, offset int, cant int) Pagina {
	return elementos[offset : offset+cant]
}

func cotizacionesToAny(cotizaciones []types.CotizacionOutputDTO) []any {
	elementos := make([]any, len(cotizaciones))
	for i, cotizacion := range cotizaciones {
		elementos[i] = cotizacion
	}
	return elementos
}
