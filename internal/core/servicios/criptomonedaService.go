package servicios

import (
	"errors"
	"fmt"
	"time"

	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/ports"
)

type MonedaServicio struct {
	repo      ports.RepositorioMonedas
	cotizador ports.Cotizador
}

func CrearServicioMoneda(r ports.RepositorioMonedas, c ports.Cotizador) MonedaServicio {
	return MonedaServicio{
		repo:      r,
		cotizador: c,
	}
}

func (s MonedaServicio) AltaMoneda(moneda, simbolo string) (int, error) {
	existe, err := s.cotizador.ExisteMoneda(simbolo)
	if err != nil {
		return 0, fmt.Errorf("error al verificar si la moneda existe: %w", err)
	}

	if !existe {
		return 0, fmt.Errorf("el simbolo %s no existe en el cotizador", simbolo)
	}

	return s.repo.AltaMoneda(domain.CrearMoneda(moneda, simbolo))
}

func (s MonedaServicio) AltaCotizacion(cotizacion domain.Cotizacion) (int, error) {
	return s.repo.AltaCotizacion(cotizacion)
}

func (s MonedaServicio) AltaCotizaciones(api string) error {
	err := s.cotizarTodas(api)
	if err != nil {
		return err
	}

	return nil
}

func (s MonedaServicio) cotizarTodas(api string) error {

	monedas, err := s.repo.BuscarTodos()
	if err != nil {
		return err
	}

	ch := make(chan error, len(monedas))
	for _, moneda := range monedas {
		go s.cotizar(api, moneda, ch)
	}
	defer close(ch)

	var errores []error
	for i := 0; i < len(monedas); i++ {
		err := <-ch
		if err != nil {
			errores = append(errores, err)
		}
	}

	return errors.Join(errores...)
}

func (s MonedaServicio) cotizar(api string, moneda domain.Criptomoneda, c chan error) {
	valor, err := s.cotizador.Cotizar(api, moneda.Simbolo)
	if err != nil {
		c <- err
		return
	}

	cotizacion := domain.Cotizacion{
		Valor:  valor,
		Time:   time.Now(),
		Moneda: moneda,
		Api:    api,
	}

	_, err = s.repo.AltaCotizacion(cotizacion)
	if err != nil {
		c <- err
		return
	}

	c <- nil
}

func (s MonedaServicio) BuscarPorId(id int) (ports.MonedaDTOOutput, error) {
	moneda, err := s.repo.BuscarPorId(id)
	if err != nil {
		return ports.MonedaDTOOutput{}, fmt.Errorf("no se encontro por id %d: %w", id, err)
	}

	return ports.MonedaDTOOutput{
		Id:           moneda.ID,
		NombreMoneda: moneda.Nombre,
		Simbolo:      moneda.Simbolo,
	}, nil
}

func (s MonedaServicio) BuscarTodos() ([]ports.MonedaDTOOutput, error) {
	monedas, err := s.repo.BuscarTodos()
	if err != nil {
		return nil, fmt.Errorf("no se pudieron obtener las monedas del repositorio: %w", err)
	}

	monedasDTOs := make([]ports.MonedaDTOOutput, len(monedas))
	for i, moneda := range monedas {
		monedasDTOs[i] = ports.MonedaDTOOutput{
			Id:           moneda.ID,
			NombreMoneda: moneda.Nombre,
			Simbolo:      moneda.Simbolo,
		}
	}

	return monedasDTOs, nil
}

func (s MonedaServicio) Cotizaciones(p ports.ParamCotizaciones) ([]ports.CotizacionOutputDTO, error) {
	cotizaciones, err := s.repo.Cotizaciones(p)
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener la cotizaciones del repositorio: %w", err)
	}

	cotizacionesDTOs := make([]ports.CotizacionOutputDTO, len(cotizaciones))
	for i, cotizacion := range cotizaciones {
		cotizacionesDTOs[i] = ports.CotizacionOutputDTO{
			NombreMoneda: cotizacion.Moneda.Nombre,
			Valor:        cotizacion.Valor,
			Fecha:        cotizacion.Time.Format("2006-01-02 15:04:05"),
			Simbolo:      cotizacion.Moneda.Simbolo,
			Api:          cotizacion.Api,
		}
	}
	return cotizacionesDTOs, nil
}

// func (s MonedaServicio) AltaCotizaciones() error {
// 	monedas, err := s.repo.BuscarTodos()
// 	if err != nil {
// 		return err
// 	}

// 	var simbolos []string
// 	for _, moneda := range monedas {
// 		simbolos = append(simbolos, moneda.Simbolo)
// 	}

// 	var errores []error
// 	cotizaciones, err := s.cotizador.CotizarConcurrente(simbolos)
// 	if err != nil {
// 		errores = append(errores, err)
// 	}

// 	if err := s.repo.InsertarCotizacionesSegunSimbolo(cotizaciones); err != nil {
// 		errores = append(errores, err)
// 	}

// 	return errors.Join(errores...)
// }
