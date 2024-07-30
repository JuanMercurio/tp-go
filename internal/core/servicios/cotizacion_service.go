package servicios

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/ports"
	"github.com/juanmercurio/tp-go/internal/ports/types"
)

type ServicioCotizacion struct {
	rc          ports.RepositorioCotizaciones
	rm          ports.RepositorioMonedas
	ru          ports.RepositorioUsuarios
	cotizadores map[string]ports.Cotizador
}

func CrearServicioCotizacion(
	rc ports.RepositorioCotizaciones,
	rm ports.RepositorioMonedas,
	ru ports.RepositorioUsuarios,
	cotizadores ...ports.Cotizador,
) *ServicioCotizacion {
	mapCotizadores := make(map[string]ports.Cotizador)
	for _, cotizador := range cotizadores {
		mapCotizadores[cotizador.GetNombre()] = cotizador
	}

	return &ServicioCotizacion{
		rc: rc,
		rm: rm,
		ru: ru,
	}
}

func (s ServicioCotizacion) AltaCotizacion(cotizacion domain.Cotizacion) (int, error) {
	return s.rc.AltaCotizacion(cotizacion)
}

func (s ServicioCotizacion) AltaCotizacionManual(usuario int, cotizacion domain.Cotizacion) (int, error) {
	return s.rc.AltaCotizacionManual(usuario, cotizacion)
}

func (s ServicioCotizacion) AltaCotizaciones(api string) error {

	monedas, err := s.rm.BuscarTodos()
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

func (s ServicioCotizacion) BajaCotizacionManual(id int) error {
	return s.rc.BajaCotizacionManual(id)
}

func (s ServicioCotizacion) PutCotizacion(username string, idCotizacion int, cotizacionNueva types.CotizacionPut) error {

	idUsuario, err := s.ru.IdDeUsername(username)
	if err != nil {
		return err
	}

	manual, err := s.rc.EsCotizacionManual(idCotizacion)
	if err != nil {
		return err
	}

	if !manual {
		return errors.New("no se puede modificar una cotizacion no manual")
	}

	cotizacionVieja, err := s.rc.CotizacionPorId(idCotizacion)
	if err != nil {
		return err
	}

	cambios := make(map[string]any)

	if cotizacionVieja.Moneda.Simbolo != cotizacionNueva.Simbolo {
		id, err := s.rm.IdDeSimbolo(cotizacionNueva.Simbolo)
		if err != nil {
			return fmt.Errorf("haciendo mapa de put: %w", err)
		}

		cambios["id_criptomoneda"] = id
	}

	fechaNueva, err := time.Parse("2006-01-02 15:04:05", cotizacionNueva.Fecha)
	if err != nil {
		return err
	}

	if cotizacionVieja.Time != fechaNueva {
		cambios["fecha"] = cotizacionNueva.Fecha
	}

	if cotizacionVieja.Valor != cotizacionNueva.Valor {
		cambios["valor"] = cotizacionNueva.Valor
	}

	return s.rc.ActualizarCotizacionMap(idUsuario, idCotizacion, cambios)
}

func (s ServicioCotizacion) Cotizaciones(p types.Filter) (int, []types.CotizacionOutputDTO, error) {
	cantidad, cotizaciones, err := s.rc.Cotizaciones(p)
	if err != nil {
		return 0, nil, fmt.Errorf("no se pudo obtener la cotizaciones del repositorio: %w", err)
	}

	cotizacionesDTOs := make([]types.CotizacionOutputDTO, len(cotizaciones))
	for i, cotizacion := range cotizaciones {
		cotizacionesDTOs[i] = types.CotizacionOutputDTO{
			NombreMoneda: cotizacion.Moneda.Nombre,
			Valor:        cotizacion.Valor,
			Fecha:        cotizacion.Time.Format("2006-01-02 15:04:05"),
			Simbolo:      cotizacion.Moneda.Simbolo,
			Api:          cotizacion.Api,
		}
	}
	return cantidad, cotizacionesDTOs, nil
}

func (s ServicioCotizacion) CotizarNuevaMoneda(id int, simbolo string) error {
	//todo implementar
	return nil
}

func (s ServicioCotizacion) Resumen(filtros types.Filter) (types.Resumen, error) {
	monedasString, fechaString, err := s.rc.Resumen(filtros)
	if err != nil {
		return types.Resumen{}, err
	}


	f, err := stringAFechasResumen(fechaString)
	if err != nil {
		return types.Resumen{}, err
	}

	m, err := stringAMapMonedasResumen(monedasString)
	if err != nil {
		return types.Resumen{}, err
	}

	r := make(types.Resumen)
	r["fechas"] = f
	r["monedas"] = m

	return r, nil
}

func (s ServicioCotizacion) CotizarManualmente(usuario int, simbolo string, fecha time.Time, precio float64) (int, error) {

	ids, err := s.rm.IdsDeSimbolos(strings.Split(simbolo, " "))
	if err != nil {
		return 0, err
	}

	moneda, err := s.rm.BuscarPorId(ids[0])
	if err != nil {
		return 0, err
	}

	cotizacion := domain.Cotizacion{
		Valor:  precio,
		Time:   fecha,
		Moneda: moneda,
		Api:    "manual",
	}

	id, err := s.AltaCotizacionManual(usuario, cotizacion)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s ServicioCotizacion) cotizarRoutine(api string,
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
func (s ServicioCotizacion) cotizarMonedaYPersistir(api string, moneda domain.Criptomoneda) (float64, error) {

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

	_, err = s.rc.AltaCotizacion(cotizacion)
	if err != nil {
		return 0, err
	}

	return valor, nil
}

func (s ServicioCotizacion) CotizacionPorId(idCotizacion int) (domain.Cotizacion, error) {
	cotizacion, err := s.rc.CotizacionPorId(idCotizacion)
	if err != nil {
		return domain.Cotizacion{}, err
	}

	cotizacion.Moneda, err = s.rm.BuscarPorId(cotizacion.Moneda.ID)
	if err != nil {
		return domain.Cotizacion{}, err
	}

	return cotizacion, nil
}
