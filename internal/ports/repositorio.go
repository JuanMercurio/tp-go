package ports

import (
	"github.com/juanmercurio/tp-go/internal/core/domain"
)

type RepositorioMonedas interface {
	AltaMoneda(domain.Criptomoneda) (int, error)
	AltaCotizacion(domain.Cotizacion) (int, error)
	AltaCotizaciones([]domain.Cotizacion) error
	BuscarPorId(int) (domain.Criptomoneda, error)
	BuscarTodos() ([]domain.Criptomoneda, error)
	Cotizaciones(Filter) (int, []domain.Cotizacion, error)
	AltaUsuario(domain.Usuario) (int, error)
	BajaUsuario(int) error
	Simbolos() []string
}
