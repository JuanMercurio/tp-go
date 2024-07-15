package repos

import (
	"database/sql"
	"fmt"
	"strings"

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

func (r RepositorioMoneda) Cotizaciones(parametros ports.ParamCotizaciones) ([]domain.Cotizacion, error) {
	SQLParams := aSQL(parametros)
	q := queryBaseCotizaciones(parametros)

	rows, err := r.db.Query(q, SQLParams.Monedas, SQLParams.FechaInicio, SQLParams.FechaFinal, SQLParams.Cant, SQLParams.Offset)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando la query en cotizaciones: %w", err)
	}
	defer rows.Close()

	cotizaciones := r.extraerCotizaciones(rows)

	return cotizaciones, nil
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
		fmt.Println("Error al preparar la sentencia:", err)
		return err
	}
	defer stmt.Close()

	tx, err := r.db.Begin()
	if err != nil {
		fmt.Println("Error al iniciar transacción:", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			fmt.Println("Transacción revertida debido a un error:", err)
		} else {
			err = tx.Commit()
			if err != nil {
				fmt.Println("Error al hacer commit de la transacción:", err)
			}
		}
	}()

	for _, row := range data {
		_, err := stmt.Exec(columns(row)...)
		if err != nil {
			fmt.Println("Error al ejecutar la sentencia:", err)
			return err
		}
	}

	return nil
}
