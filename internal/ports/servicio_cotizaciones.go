package ports

import (
	"time"

	"github.com/juanmercurio/tp-go/internal/ports/types"
)

type ServicioCotizacion interface {
	Cotizaciones(types.Filter) (int, []types.CotizacionOutputDTO, error)
	AltaCotizaciones(api string) error
	CotizarNuevaMoneda(id int, simbolo string) error
	Resumen(filtros types.Filter) (types.Resumen, error)
	CotizarManualmente(usuario int, simbolo string, fecha time.Time, precio float64) (int, error)
	BajaCotizacionManual(id int) error
	PutCotizacion(username string, idCotizacion int, cotizacion types.CotizacionPut) error
}
