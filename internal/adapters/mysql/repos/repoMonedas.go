package repos

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/ports"
)

type RepositorioMoneda struct {
	db *sql.DB
}

func CrearRepositorioMoneda(db *sql.DB) RepositorioMoneda {
	return RepositorioMoneda{
		db: db,
	}
}

func (r RepositorioMoneda) AltaMoneda(moneda domain.Criptomoneda) (int, error) {
	query := "INSERT INTO go.criptomoneda (nombre) VALUES (?)"
	return r.darDeAltaEntidad(query, moneda.Nombre)
}

func (r RepositorioMoneda) AltaCotizacion(cotizacion domain.Cotizacion) (int, error) {
	query := "INSERT INTO go.cotizacion (id_criptomoneda, fecha, valor) VALUES (?, ?, ?)"
	return r.darDeAltaEntidad(query, cotizacion.Moneda.ID, cotizacion.Time, cotizacion.Valor)
}

func (r RepositorioMoneda) darDeAltaEntidad(query string, params ...any) (int, error) {
	result, err := r.db.Exec(query, params...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r RepositorioMoneda) BuscarPorId(id int) (domain.Criptomoneda, error) {
	query := "SELECT id, nombre FROM go.criptomoneda WHERE id = ?"

	var moneda domain.Criptomoneda

	err := r.db.QueryRow(query, id).Scan(&moneda.ID, &moneda.Nombre)
	if err != nil {
		return domain.Criptomoneda{}, err
	}

	return moneda, nil
}

func (r RepositorioMoneda) BuscarTodos() ([]domain.Criptomoneda, error) {
	rows, err := r.db.Query("SELECT * FROM go.criptomoneda")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monedas []domain.Criptomoneda

	for rows.Next() {
		moneda := domain.Criptomoneda{}
		if err := rows.Scan(&moneda.ID, &moneda.Nombre); err != nil {
			return nil, err
		}
		monedas = append(monedas, moneda)
	}

	return monedas, nil
}

func (r RepositorioMoneda) obtenerCotizaciones(rows *sql.Rows) ([]domain.Cotizacion, error) {

	var cotizaciones []domain.Cotizacion
	for rows.Next() {

		var cotizacion domain.Cotizacion
		var id_criptomoneda int
		var tiempoString string

		if err := rows.Scan(&cotizacion.ID, &id_criptomoneda, &tiempoString, &cotizacion.Valor); err != nil {
			return nil, err
		}

		moneda, err := r.BuscarPorId(id_criptomoneda)
		if err != nil {
			return nil, err
		}

		cotizacion.Time, err = time.Parse("2006-01-02 15:04:05", tiempoString)
		if err != nil {
			return nil, err
		}

		cotizacion.Moneda = moneda
		cotizaciones = append(cotizaciones, cotizacion)
	}

	if len(cotizaciones) == 0 {
		return nil, errors.New("no existen cotizaciones")
	}

	return cotizaciones, nil
}

// struct?
func (r RepositorioMoneda) Cotizaciones(p ports.Parametros) ([]domain.Cotizacion, error) {
	monedas := p.Monedas
	fechaInicio := p.FechaInicial
	fechaFinal := p.FechaFinal
	cant := p.CantPaginas * p.TamPaginas

	query := fmt.Sprintf(
		`SELECT * FROM go.cotizacion
		 WHERE id_criptomoneda IN (%s) 
		 AND fecha BETWEEN ? AND ?
		 ORDER BY fecha DESC LIMIT ?`,
		strings.Join(monedas, ","))

	rows, err := r.db.Query(query, fechaInicio, fechaFinal, cant)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando la query en cotizaciones: %w", err)
	}
	defer rows.Close()

	cotizaciones := r.extraerCotizaciones(rows)
	// if len(cotizaciones) == 0 {
	// 	return nil, errors.New("no existen cotizaciones para estos filtros")
	// }

	return cotizaciones, nil
}

func (r RepositorioMoneda) extraerCotizaciones(rows *sql.Rows) []domain.Cotizacion {
	var cotizaciones []domain.Cotizacion
	for rows.Next() {
		var cotizacion domain.Cotizacion
		var id_moneda int
		var tiempoString string
		if err := rows.Scan(&cotizacion.ID, &id_moneda, &tiempoString, &cotizacion.Valor); err != nil {
			log.Fatal(err)
		}

		// TODO no ir siempre a la base esto puede ser mas performante si memoizamos
		moneda, _ := r.BuscarPorId(id_moneda)
		cotizacion.Time, _ = time.Parse("2006-01-02 15:04:05", tiempoString)
		cotizacion.Moneda = moneda

		cotizaciones = append(cotizaciones, cotizacion)
	}

	return cotizaciones
}

func (r RepositorioMoneda) BuscarCotizacionMoneda(
	moneda int,
	fechaInicial time.Time,
	fechaFinal time.Time,
	cant int,
	offset int,
) ([]domain.Cotizacion, error) {

	query := "SELECT * FROM go.cotizacion WHERE id_criptomoneda = ? AND fecha BETWEEN ? AND ? LIMIT ? OFFSET ?"

	rows, err := r.db.Query(query, moneda, fechaInicial, fechaFinal, cant, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.obtenerCotizaciones(rows)
}
