package ports

type ServicioMonedas interface {
	BuscarTodos() ([]MonedaDTOOutput, error)
	AltaMoneda(nombreMoneda, simbolo string) (int, error)
	Cotizaciones(ParamCotizaciones) ([]CotizacionOutputDTO, error)
	AltaCotizaciones(api string) error
}
