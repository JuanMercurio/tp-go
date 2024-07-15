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

func (s MonedaServicio) AltaMoneda(moneda string) (int, error) {
	return s.repo.AltaMoneda(domain.CrearMoneda(moneda))
}

func (s MonedaServicio) AltaCotizacion(cotizacion domain.Cotizacion) (int, error) {
	return s.repo.AltaCotizacion(cotizacion)
}

func (s MonedaServicio) AltaCotizaciones() ([]domain.Cotizacion, error) {
	cotizaciones, err := s.BuscarCotizacionesEnAPIs()
	if err != nil {
		return nil, err
	}

	//TODO deberia retornar con el id todo?
	err = s.repo.AltaCotizaciones(cotizaciones)
	if err != nil {
		return nil, err
	}

	return cotizaciones, err
}

func (s MonedaServicio) BuscarCotizacionesEnAPIs() ([]domain.Cotizacion, error) {
	moneda := domain.Criptomoneda{
		Nombre: "Bitcoin",
		ID:     1,
	}

	cotizaciones := []domain.Cotizacion{
		{
			Moneda: moneda,
			Valor:  1.1,
			Time:   time.Now(),
		},

		{
			Moneda: moneda,
			Valor:  234,
			Time:   time.Now(),
		},
	}
	return cotizaciones, nil
}

func (s MonedaServicio) BuscarPorId(id int) (domain.Criptomoneda, error) {
	return s.repo.BuscarPorId(id)
}

func (s MonedaServicio) BuscarTodos() ([]domain.Criptomoneda, error) {
	return s.repo.BuscarTodos()
}

func (s MonedaServicio) Cotizaciones(p ports.ParamCotizaciones) ([]domain.Cotizacion, error) {
	return s.repo.Cotizaciones(p)
}
