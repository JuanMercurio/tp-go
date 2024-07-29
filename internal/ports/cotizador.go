//go:generate  /home/feb/go/bin/mockgen  --destination=./mock/cotizador.go github.com/juanmercurio/tp-go/internal/ports Cotizador
package ports

type Cotizacion struct {
	Valor   float64
	Simbolo string
	Err     error
}

type Cotizador interface {
	ExisteMoneda(simbolo string) (bool, error)
	Cotizar(simbolos string) (float64, error)
	GetNombre() string
}
