package ports

import "time"

type ServicioCotizacion interface {
	Cotizaciones(Filter) (int, []CotizacionOutputDTO, error)
	AltaCotizaciones(api string) error
	CotizarNuevaMoneda(id int, simbolo string) error
	Resumen(filtros Filter) (Resumen, error)
	CotizarManualmente(usuario int, simbolo string, fecha time.Time, precio float64) error
	BajaCotizacionManual(id int) error
	PutCotizacion(idUsuario int, idCotizacion int, cotizacion CotizacionPut) error
}
