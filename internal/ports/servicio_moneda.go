package ports

type ServicioMonedas interface {
	AltaMoneda(nombreMoneda, simbolo string) (int, error)
	BuscarTodos() ([]MonedaOutputDTO, error)
	SimbolosValido(simbolos []string) error
	SimboloValido(simbolos string) error
}
