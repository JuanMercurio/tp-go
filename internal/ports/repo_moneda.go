//go:generate  /home/feb/go/bin/mockgen  --destination=./mock/repo_monedas.go github.com/juanmercurio/tp-go/internal/ports RepositorioMonedas
package ports

import "github.com/juanmercurio/tp-go/internal/core/domain"

type RepositorioMonedas interface {
	AltaMoneda(domain.Criptomoneda) (int, error)
	IdsDeSimbolos(simbolos []string) ([]int, error)
	IdDeSimbolo(simbolo string) (int, error)
	SimbolosValido(simbolos []string) error
	SimboloValido(simbolos string) error
	BuscarPorId(int) (domain.Criptomoneda, error)
	BuscarTodos() ([]domain.Criptomoneda, error)
}
