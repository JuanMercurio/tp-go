package ports

type ServicioMonedas interface {
	AltaMoneda(nombreMoneda, simbolo string) (int, error)
	BuscarTodos() ([]MonedaOutputDTO, error)
	Cotizaciones(Filter) (int, []CotizacionOutputDTO, error)
	AltaCotizaciones(api string) error
	CotizarNuevaMoneda(id int, simbolo string) error
	MonedasDeUsuario(id int) ([]MonedaOutputDTO, error)
}
