package servicios

import (
	"time"

	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/ports"
)

type MonedaServicio struct {
	repo ports.RepositorioMonedas
}

//TODO: no se esta teniendo en cuenta dtos

func CrearServicioMoneda(r ports.RepositorioMonedas) MonedaServicio {
	return MonedaServicio{
		repo: r,
	}
}

func (s MonedaServicio) AltaMoneda(moneda domain.Criptomoneda) (int, error) {
	return s.repo.AltaMoneda(moneda)
}

func (s MonedaServicio) AltaCotizacion(cotizacion domain.Cotizacion) (int, error) {
	return s.repo.AltaCotizacion(cotizacion)
}

func (s MonedaServicio) BuscarPorId(id int) (domain.Criptomoneda, error) {
	return s.repo.BuscarPorId(id)
}

func (s MonedaServicio) BuscarTodos() ([]domain.Criptomoneda, error) {
	return s.repo.BuscarTodos()
}

func (s MonedaServicio) BuscarCotizacionMoneda(
	monedaID int,
	fechaInicial time.Time,
	fechaFinal time.Time,
	cant int,
	offset int,
) ([]domain.Cotizacion, error) {
	return s.repo.BuscarCotizacionMoneda(monedaID, fechaInicial, fechaFinal, cant, offset)
}

func (s MonedaServicio) Cotizaciones(p ports.Parametros) ([]domain.Cotizacion, error) {
	return s.repo.Cotizaciones(p)
}
