package ports

type ServicioMonedas interface {
	BuscarTodos() ([]MonedaOutputDTO, error)
	AltaMoneda(nombreMoneda, simbolo string) (int, error)
	Cotizaciones(ParamCotizaciones) ([]CotizacionOutputDTO, error)
	AltaCotizaciones(api string) error
	CotizarNuevaMoneda(simbolo string) error
}
