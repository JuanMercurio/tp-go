//go:generate  /home/feb/go/bin/mockgen  --destination=./mock/repo_monedas.go github.com/juanmercurio/tp-go/internal/ports RepositorioMonedas
package ports

import "github.com/juanmercurio/tp-go/internal/core/domain"

type RepositorioMonedas interface {
	AltaMoneda(domain.Criptomoneda) (int, error)
	AltaCotizacion(domain.Cotizacion) (int, error)
	BuscarPorId(int) (domain.Criptomoneda, error)
	BuscarTodos() ([]domain.Criptomoneda, error)
	Cotizaciones(Filter) (int, []domain.Cotizacion, error)
	Resumen(Filter) (string, string, error)
	SimbolosValido(simbolos []string) error
	SimboloValido(simbolos string) error
	IdsDeSimbolos(simbolos []string) ([]int, error)
	BajaCotizacionManual(id int) error
	ActualizarCotizacionMap(idUsuario int, idCotizacion int, cambios map[string]any) error
}
