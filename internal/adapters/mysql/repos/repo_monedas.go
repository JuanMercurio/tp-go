package repos

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

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

func (r RepositorioMoneda) Cotizaciones(parametros ports.ParamCotizaciones) ([]domain.Cotizacion, error) {
	q := QueryBaseCotizaciones(parametros).toString()
	rows, err := r.db.Query(q, parametros.FechaInicial, parametros.FechaFinal, parametros.TamPaginas*parametros.CantPaginas, parametros.PaginaInicial)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando la query en cotizaciones: %w", err)
	}
	defer rows.Close()

	cotizaciones := r.extraerCotizaciones(rows)

	return cotizaciones, nil
}

// esto no va

// polemico porque hace este insert por ejemplo
// INSERT INTO go.cotizacion (id_criptomoneda,fecha,  valor)
// VALUES ((SELECT id FROM go.criptomoneda where simbolo = "BTC"),"2024-07-15 23:34:56", 64581.195000)
func (r RepositorioMoneda) InsertarCotizacionesSegunSimbolo(c []ports.Cotizacion) error {
	queryBase := "INSERT INTO go.cotizacion (id_criptomoneda,fecha,  valor) VALUES "
	for _, cotizacion := range c {
		queryBase += fmt.Sprintf(`((SELECT id FROM go.criptomoneda where simbolo = "%s"),"%s", %f),`,
			cotizacion.Simbolo,
			time.Now().Format("2006-01-02 15:04:05"),
			cotizacion.Valor)

	}
	query := strings.TrimSuffix(queryBase, ",")

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r RepositorioMoneda) AltaCotizaciones(cotizaciones []domain.Cotizacion) error {
	converter := func(c any) []any {
		var columnas []any
		cotizacion := c.(domain.Cotizacion)
		columnas = append(columnas, any(cotizacion.Moneda.ID))
		columnas = append(columnas, any(cotizacion.Time))
		columnas = append(columnas, any(cotizacion.Valor))
		return columnas
	}

	cotizacionesAny := make([]any, len(cotizaciones))
	for i, c := range cotizaciones {
		cotizacionesAny[i] = any(c)
	}

	err := r.AltaMasivaCustomColumns("go.cotizacion", "id_criptomoneda, fecha, valor", cotizacionesAny, converter)
	if err != nil {
		return fmt.Errorf("error en el alta de la cotizacion: %w", err)
	}

	return nil
}

func (r RepositorioMoneda) AltaMasivaCustomColumns(tabla string, columnas string, data []any, columns func(any) []any) error {
	var valores []string
	for range strings.Split(columnas, " ") {
		valores = append(valores, "?")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tabla, columnas, strings.Join(valores, ","))
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
			}
		}
	}()

	for _, row := range data {
		_, err := stmt.Exec(columns(row)...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r RepositorioMoneda) Simbolos() []string {
	query := "SELECT simbolo FROM go.criptomoneda"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var simbolos []string

	for rows.Next() {
		var simbolo string
		if err := rows.Scan(&simbolo); err != nil {
			return nil
		}
		simbolos = append(simbolos, simbolo)
	}

	return simbolos
}
