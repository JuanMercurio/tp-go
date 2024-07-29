package ports

import "github.com/juanmercurio/tp-go/internal/ports/types"

type ServicioMonedas interface {
	AltaMoneda(nombreMoneda, simbolo string) (int, error)
	BuscarTodos() ([]types.MonedaOutputDTO, error)
	SimbolosValido(simbolos []string) error
	SimboloValido(simbolos string) error
}
