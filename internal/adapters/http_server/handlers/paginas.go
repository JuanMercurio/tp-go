package handlers

import (
	"math"

	"github.com/juanmercurio/tp-go/internal/core/domain"
)

type Pagina []any

func crearPaginas(tamPaginas int, cotizaciones []domain.Cotizacion) []Pagina {
	cotizacionesAny := deCotizacionAAny(cotizaciones)
	var paginas []Pagina
	restantes := len(cotizaciones)
	cantidadPaginasFinal := paginasNecesarias(tamPaginas, len(cotizaciones))

	for i := 0; i < cantidadPaginasFinal; i++ {
		offset := tamPaginas * i
		if restantes < tamPaginas {
			tamPaginas = restantes
		}
		p := obtenerPaginas(cotizacionesAny, offset, tamPaginas)
		paginas = append(paginas, p)
		restantes -= tamPaginas
	}

	return paginas
}
func paginasNecesarias(tamPaginas int, cantCotizaciones int) int {
	if cantCotizaciones < tamPaginas {
		return 1
	}
	return int(math.Ceil(float64(cantCotizaciones) / float64(tamPaginas)))
}

func obtenerPaginas(cotizaciones []any, offset int, cant int) Pagina {
	return cotizaciones[offset : offset+cant]
}

func deCotizacionAAny(cotizaciones []domain.Cotizacion) []any {
	s := make([]any, len(cotizaciones))
	for i, cotizacion := range cotizaciones {
		s[i] = cotizacion
	}
	return s

}
