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
	// CotizarConcurrente(simbolos []string) ([]Cotizacion, error)
}
