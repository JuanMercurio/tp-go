//go:generate  /home/feb/go/bin/mockgen  --destination=./mock/repo_cotizaciones.go github.com/juanmercurio/tp-go/internal/ports RepositorioCotizaciones
package ports

import (
	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/ports/types"
)

type RepositorioCotizaciones interface {
	AltaCotizacion(cotizacion domain.Cotizacion) (int, error)

	Cotizaciones(types.Filter) (int, []domain.Cotizacion, error)
	Resumen(types.Filter) (string, string, error)

	AltaCotizacionManual(id int, cotizacion domain.Cotizacion) (int, error)
	BajaCotizacionManual(id int) error
	ActualizarCotizacionMap(idUsuario int, idCotizacion int, cambios map[string]any) error

	CotizacionPorId(id int) (domain.Cotizacion, error)
	EsCotizacionManual(id int) (bool, error)
}
