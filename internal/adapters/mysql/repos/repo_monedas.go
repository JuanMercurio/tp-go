package repos

import (
	"database/sql"
	"fmt"

	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/ports"
)

type RepositorioMoneda struct {
	db *sql.DB
}

func CrearRepositorioMoneda(db *sql.DB) *RepositorioMoneda {
	return &RepositorioMoneda{
		db: db,
	}
}

func (r RepositorioMoneda) AltaMoneda(moneda domain.Criptomoneda) (int, error) {
	query := "INSERT INTO go.criptomoneda (nombre, simbolo) VALUES (?, ?)"
	return r.darDeAltaEntidad(query, moneda.Nombre, moneda.Simbolo)
}

func (r RepositorioMoneda) AltaCotizacion(cotizacion domain.Cotizacion) (int, error) {
	query := "INSERT INTO go.cotizacion (id_criptomoneda, fecha, valor, api) VALUES (?, ?, ?, ?)"
	return r.darDeAltaEntidad(query, cotizacion.Moneda.ID, cotizacion.Time, cotizacion.Valor, cotizacion.Api)
}

func (r RepositorioMoneda) darDeAltaEntidad(query string, params ...any) (int, error) {
	result, err := r.db.Exec(query, params...)
	if err != nil {
		return 0, fmt.Errorf("error al ejecutar el query en la base: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al buscar el id de la moneda insertada: %w", err)
	}

	return int(id), nil
}

func (r RepositorioMoneda) BuscarPorId(id int) (domain.Criptomoneda, error) {
	query := "SELECT id, nombre, simbolo FROM go.criptomoneda WHERE id = ?"

	var moneda domain.Criptomoneda

	err := r.db.QueryRow(query, id).Scan(&moneda.ID, &moneda.Nombre, &moneda.Simbolo)
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
		if err := rows.Scan(&moneda.ID, &moneda.Nombre, &moneda.Simbolo); err != nil {
			return nil, err
		}
		monedas = append(monedas, moneda)
	}

	return monedas, nil
}

func (r RepositorioMoneda) Cotizaciones(parametros ports.Filter) (int, []domain.Cotizacion, error) {

	//TODO modularizar este codigo
	q := QueryBaseCotizaciones(parametros)

	q.Select = "SELECT COUNT(*)"
	var cantidad int

	err := r.db.QueryRow(q.toString()).Scan(&cantidad)
	if err != nil {
		return 0, nil, fmt.Errorf("error ejecutando la query en cotizaciones: %w", err)
	}

	q.Select = "SELECT cotizacion.id_criptomoneda, fecha, valor, api"
	q.AddLimit(parametros.TamPaginas * parametros.CantPaginas)
	q.AddOffset(parametros.PaginaInicial)

	rows, err := r.db.Query(q.toString())
	if err != nil {
		return 0, nil, fmt.Errorf("error ejecutando la query en cotizaciones: %w", err)
	}
	defer rows.Close()

	cotizaciones := r.extraerCotizaciones(rows)

	return cantidad, cotizaciones, nil
}

func (r RepositorioMoneda) MonedasDeUsuario(id int) ([]domain.Criptomoneda, error) {

	var builder queryBuilder
	builder.Select = "SELECT criptomoneda.id, criptomoneda.nombre, criptomoneda.simbolo"
	builder.From = "FROM criptomoneda"
	builder.AddJoin("usuario_criptomoneda", "criptomoneda.id = usuario_criptomoneda.id_criptomoneda")
	builder.AddWhere("usuario_criptomoneda.id_usuario = ?")

	query := builder.toString()

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var moneda domain.Criptomoneda
	var monedas []domain.Criptomoneda
	for rows.Next() {
		if err := rows.Scan(&moneda.ID, &moneda.Nombre, &moneda.Simbolo); err != nil {
			return nil, err
		}
		monedas = append(monedas, moneda)
	}

	return monedas, nil
}
