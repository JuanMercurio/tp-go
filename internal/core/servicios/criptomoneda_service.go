package servicios

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
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

	monedas, err := s.repo.BuscarTodos()
	if err != nil {
		return err
	}

	//TODO la cantidad de go routines esta hardcodeada
	cantRoutines := 10

	chErrors := make(chan error)
	chMonedas := make(chan domain.Criptomoneda)

	var wg sync.WaitGroup
	wg.Add(cantRoutines)

	for i := 0; i < cantRoutines; i++ {
		go s.cotizarRoutine(api, chMonedas, chErrors, &wg)
	}

	for _, moneda := range monedas {
		chMonedas <- moneda
	}
	close(chMonedas)

	wg.Wait()
	close(chErrors)

	var errores []error
	for i := 0; i < len(monedas); i++ {
		err := <-chErrors
		if err != nil {
			errores = append(errores, err)
		}
	}

	return errors.Join(errores...)
}

func (s MonedaServicio) cotizarRoutine(api string,
	monedasCanal chan domain.Criptomoneda,
	errorCanal chan error,
	wg *sync.WaitGroup) {

	for moneda := range monedasCanal {
		if _, err := s.cotizarMonedaYPersistir(api, moneda); err != nil {
			errorCanal <- err
		}
	}

	wg.Done()
}

// TODO seria mejor si retorna la cotizacion entera
func (s MonedaServicio) cotizarMonedaYPersistir(api string, moneda domain.Criptomoneda) (float64, error) {

	cotizador, ok := s.cotizadores[api]
	if !ok {
		return 0, fmt.Errorf("no existe el cotizador " + api)
	}

	valor, err := cotizador.Cotizar(moneda.Simbolo)
	if err != nil {
		return 0, err
	}

	cotizacion := domain.Cotizacion{
		Valor:  valor,
		Time:   time.Now(),
		Moneda: moneda,
		Api:    api,
	}

	_, err = s.repo.AltaCotizacion(cotizacion)
	if err != nil {
		return 0, err
	}

	return valor, nil
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

func (s MonedaServicio) Cotizaciones(p ports.Filter) (int, []ports.CotizacionOutputDTO, error) {
	cantidad, cotizaciones, err := s.repo.Cotizaciones(p)
	if err != nil {
		return 0, nil, fmt.Errorf("no se pudo obtener la cotizaciones del repositorio: %w", err)
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
	return cantidad, cotizacionesDTOs, nil
}

func (s MonedaServicio) CotizarNuevaMoneda(id int, simbolo string) error {
	//todo implementar
	return nil
}

func (s MonedaServicio) Resumen(filtros ports.Filter) (ports.Resumen, error) {
	monedasString, fechaString, err := s.repo.Resumen(filtros)
	if err != nil {
		return ports.Resumen{}, err
	}

	fmt.Println(fechaString)

	f, err := stringAFechasResumen(fechaString)
	if err != nil {
		return ports.Resumen{}, err
	}

	m, err := stringAMapMonedasResumen(monedasString)
	if err != nil {
		return ports.Resumen{}, err
	}

	r := make(ports.Resumen)
	r["fechas"] = f
	r["monedas"] = m

	return r, nil
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

func (s MonedaServicio) CotizarManualmente(simbolo string, fecha time.Time, precio float64) error {

	ids, err := s.repo.IdsDeSimbolos(strings.Split(simbolo, " "))
	if err != nil {
		return err
	}

	moneda, err := s.repo.BuscarPorId(ids[0])
	if err != nil {
		return err
	}

	cotizacion := domain.Cotizacion{
		Valor:  precio,
		Time:   fecha,
		Moneda: moneda,
		Api:    "manual",
	}

	_, err = s.AltaCotizacion(cotizacion)
	if err != nil {
		return err
	}

	return nil
}

func (s MonedaServicio) BajaCotizacionManual(id int) error {
	return s.repo.BajaCotizacionManual(id)
}

func (s MonedaServicio) PatchCotizacion(idUsuario int, idCotizacion int, patchs []ports.Patch) error {

	m := make(map[string]any)
	for _, patch := range patchs {
		m[patch.Path] = patch.NuevoValor
	}

	if err := s.repo.ActualizarCotizacionMap(idUsuario, idCotizacion, m); err != nil {
		return err
	}

	return nil
}
