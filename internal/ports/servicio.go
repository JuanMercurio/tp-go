package ports

type ServicioMonedas interface {
	AltaUsuario(nombre string) (int, error)
	BajaUsuario(id int) error
	BuscarTodos() ([]MonedaOutputDTO, error)
	AltaMoneda(nombreMoneda, simbolo string) (int, error)
	Cotizaciones(Filter) (int, []CotizacionOutputDTO, error)
	AltaCotizaciones(api string) error
	CotizarNuevaMoneda(id int, simbolo string) error
}
