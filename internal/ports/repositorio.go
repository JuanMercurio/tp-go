package ports

import (
	"time"

	"github.com/juanmercurio/tp-go/internal/core/domain"
)

type RepositorioMonedas interface {
	AltaMoneda(domain.Criptomoneda) (int, error)
	AltaCotizacion(domain.Cotizacion) (int, error)
	BuscarPorId(int) (domain.Criptomoneda, error)
	BuscarTodos() ([]domain.Criptomoneda, error)
	Cotizaciones(Parametros) ([]domain.Cotizacion, error)
	BuscarCotizacionMoneda(moneda int, fechaInicial time.Time, fechaFinal time.Time, cant int, offset int) ([]domain.Cotizacion, error)
}
