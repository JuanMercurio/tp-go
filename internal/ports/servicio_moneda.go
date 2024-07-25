package ports

import "time"

type ServicioMonedas interface {
	AltaMoneda(nombreMoneda, simbolo string) (int, error)
	BuscarTodos() ([]MonedaOutputDTO, error)
	Cotizaciones(Filter) (int, []CotizacionOutputDTO, error)
	AltaCotizaciones(api string) error
	CotizarNuevaMoneda(id int, simbolo string) error
	Resumen(filtros Filter) (Resumen, error)
	SimbolosValido(simbolos []string) error
	SimboloValido(simbolos string) error
	CotizarManualmente(simbolo string, fecha time.Time, precio float64) error
	BajaCotizacionManual(id int) error
	PatchCotizacion(idUsuario int, idCotizacion int, patchs []Patch) error
}
