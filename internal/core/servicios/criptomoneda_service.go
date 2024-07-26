package servicios

import (
	"encoding/json"
	"fmt"

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

func stringAMapMonedasResumen(monedasString string) (map[string]int, error) {
	var data []string

	if err := json.Unmarshal([]byte(monedasString), &data); err != nil {
		return nil, err
	}

	m := make(map[string]int)
	for _, value := range data {
		if _, ok := m[value]; !ok {
			m[value] = 1
		} else {
			m[value]++
		}
	}

	return m, nil
}

func stringAFechasResumen(fechaString string) (map[string]any, error) {
	var fechas map[string]any
	if err := json.Unmarshal([]byte(fechaString), &fechas); err != nil {
		return nil, err
	}
	return fechas, nil
}

func (s MonedaServicio) SimboloValido(simbolo string) error {
	return s.repo.SimboloValido(simbolo)
}

func (s MonedaServicio) SimbolosValido(simbolos []string) error {
	return s.repo.SimbolosValido(simbolos)
}
