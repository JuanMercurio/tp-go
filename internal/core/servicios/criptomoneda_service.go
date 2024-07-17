package servicios

import (
	"errors"
	"fmt"
	"time"

	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/ports"
)

type MonedaServicio struct {
	repo        ports.RepositorioMonedas
	cotizadores map[string]ports.Cotizador
}

func CrearServicioMoneda(r ports.RepositorioMonedas, cotizadores ...ports.Cotizador) *MonedaServicio {
	mapCotizadores := make(map[string]ports.Cotizador)
	for _, cotizador := range cotizadores {
		mapCotizadores[cotizador.GetNombre()] = cotizador
	}
	return &MonedaServicio{
		repo:        r,
		cotizadores: mapCotizadores,
	}
}

func (s MonedaServicio) AltaMoneda(moneda, simbolo string) (int, error) {

	for nombre, cotizador := range s.cotizadores {

		fmt.Println("Estoy en el cotizador", nombre)
		existe, err := cotizador.ExisteMoneda(simbolo)
		if err != nil {
			return 0, fmt.Errorf("error al verificar si la moneda existe: %w", err)
		}

		if !existe {
			return 0, fmt.Errorf("el simbolo %s no existe en el cotizador: %s", simbolo, nombre)
		}
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
		go s.cotizarRoutine(api, moneda, ch)
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

func (s MonedaServicio) cotizarRoutine(api string, moneda domain.Criptomoneda, c chan error) {
	if err := s.CotizarMoneda(api, moneda); err != nil {
		c <- err
		return
	}
	c <- nil
}

func (s MonedaServicio) CotizarMoneda(api string, moneda domain.Criptomoneda) error {

	cotizador, ok := s.cotizadores[api]
	if !ok {
		return fmt.Errorf("no existe el cotizador " + api)
	}

	valor, err := cotizador.Cotizar(moneda.Simbolo)
	if err != nil {
		return err
	}

	cotizacion := domain.Cotizacion{
		Valor:  valor,
		Time:   time.Now(),
		Moneda: moneda,
		Api:    api,
	}

	_, err = s.repo.AltaCotizacion(cotizacion)
	if err != nil {
		return err
	}

	return nil
}

func (S MonedaServicio) CotizarNuevaMoneda(simbolo string) error {
	return nil
}

func (s MonedaServicio) BuscarPorId(id int) (ports.MonedaOutputDTO, error) {
	moneda, err := s.repo.BuscarPorId(id)
	if err != nil {
		return ports.MonedaOutputDTO{}, fmt.Errorf("no se encontro por id %d: %w", id, err)
	}

	return ports.MonedaOutputDTO{
		Id:           moneda.ID,
		NombreMoneda: moneda.Nombre,
		Simbolo:      moneda.Simbolo,
	}, nil
}

func (s MonedaServicio) BuscarTodos() ([]ports.MonedaOutputDTO, error) {
	monedas, err := s.repo.BuscarTodos()
	if err != nil {
		return nil, fmt.Errorf("no se pudieron obtener las monedas del repositorio: %w", err)
	}

	monedasDTOs := make([]ports.MonedaOutputDTO, len(monedas))
	for i, moneda := range monedas {
		monedasDTOs[i] = ports.MonedaOutputDTO{
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
