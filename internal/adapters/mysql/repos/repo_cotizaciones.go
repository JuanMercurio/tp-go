package repos

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/ports/types"
)

type RepositorioCotizaciones struct {
	db *sql.DB
}

func CrearRepositorioCotizaciones(cliente *sql.DB) *RepositorioCotizaciones {
	return &RepositorioCotizaciones{db: cliente}
}

func (r RepositorioCotizaciones) AltaCotizacion(cotizacion domain.Cotizacion) (int, error) {
	query := "INSERT INTO go.cotizacion (id_criptomoneda, fecha, valor, api) VALUES (?, ?, ?, ?)"

	result, err := r.db.Exec(query, cotizacion.Moneda.ID, cotizacion.Time, cotizacion.Valor, cotizacion.Api)
	if err != nil {
		return 0, fmt.Errorf("error al ejecutar el query en la base: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al buscar el id de la moneda insertada: %w", err)
	}

	return int(id), nil
}

func (r RepositorioCotizaciones) AltaCotizacionManual(usuario int, cotizacion domain.Cotizacion) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	cotizacion.ID, err = r.AltaCotizacion(cotizacion)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := r.Auditar(usuario, cotizacion.ID, "alta", "valor", "", cotizacion.Valor); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := r.Auditar(usuario, cotizacion.ID, "alta", "fecha", "", cotizacion.Time); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return cotizacion.ID, nil
}

func (r RepositorioCotizaciones) BajaCotizacionManual(id int) error {
	query := "DELETE FROM go.cotizacion WHERE id = ? AND api = 'manual'"
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// para devolver mejores errores
	if rowsAffected == 0 {
		query := "SELECT COUNT(*) FROM go.cotizacion WHERE id = ?"
		var count int
		err := r.db.QueryRow(query, id).Scan(&count)
		if err != nil {
			return err
		}

		if count == 0 {
			return errors.New("la cotizacion no existe")
		}

		return errors.New("la cotizacion no es manual")
	}

	return nil

}

func (r RepositorioCotizaciones) CotizacionPorId(id int) (domain.Cotizacion, error) {

	query := "SELECT id_criptomoneda, fecha, valor, api FROM go.cotizacion WHERE id = ?"
	var cotizacion domain.Cotizacion
	var idMoneda int
	var fechaString string
	if err := r.db.QueryRow(query, id).Scan(&idMoneda, &fechaString, &cotizacion.Valor, &cotizacion.Api); err != nil {
		return domain.Cotizacion{}, err
	}

	var err error
	cotizacion.Time, err = time.Parse("2006-01-02 15:04:05", fechaString)
	if err != nil {
		return domain.Cotizacion{}, err
	}

	moneda, err := r.monedaPorId(idMoneda)
	if err != nil {
		return domain.Cotizacion{}, err
	}

	cotizacion.Moneda = moneda
	cotizacion.ID = id

	return cotizacion, nil
}

func (r RepositorioCotizaciones) ActualizarCotizacionMap2(usuario int, idCotizacion int, cambios map[string]any) error {
	return nil
}

func (r RepositorioCotizaciones) ActualizarCotizacionMap(usuario int, idCotizacion int, cambios map[string]interface{}) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	valoresActuales, err := r.obtenerValoresActuales(idCotizacion)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	query := "UPDATE cotizacion SET "
	params := make([]interface{}, 0)

	for campo, valor := range cambios {
		query += campo + " = ?, "
		params = append(params, valor)
	}
	query = strings.TrimSuffix(query, ", ")

	query += fmt.Sprintf(" WHERE id = %d AND api = 'manual'", idCotizacion)

	_, err = r.db.Exec(query, params...)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	for campo, valorNuevo := range cambios {
		valorAnterior := valoresActuales[campo]
		if err := r.Auditar(usuario, idCotizacion, "actualizar", campo, valorAnterior, valorNuevo); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r RepositorioCotizaciones) obtenerValoresActuales(idCotizacion int) (map[string]interface{}, error) {
	query := `SELECT fecha, valor, id_criptomoneda FROM cotizacion WHERE id = ?`

	row := r.db.QueryRow(query, idCotizacion)

	valoresActuales := make(map[string]interface{})
	var fecha, valor, idMoneda any

	err := row.Scan(&fecha, &valor, &idMoneda)
	if err != nil {
		return nil, err
	}

	valoresActuales["fecha"] = fecha
	valoresActuales["valor"] = valor
	valoresActuales["id_criptomoneda"] = idMoneda

	return valoresActuales, nil
}

func (r RepositorioCotizaciones) Auditar(usuario int, cotizacion int, accion string, columnaAfectada string, valorViejo any, nuevoValor any) error {
	query := `
	INSERT INTO auditoria (id_usuario, id_cotizacion, accion, columna_afectada, nuevo_valor, viejo_valor ,fecha) VALUES (?,?,?,?,?,?,?)	
	`

	res, err := r.db.Exec(query, usuario, cotizacion, accion, columnaAfectada, nuevoValor, valorViejo, time.Now())
	if err != nil {
		return err
	}

	afectadas, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if afectadas == 0 {
		return errors.New("no se pudo auditar")
	}

	return nil

}

func (r RepositorioCotizaciones) Cotizaciones(parametros types.Filter) (int, []domain.Cotizacion, error) {

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

// la version facil del select seria --> select group_concat(distinct id_criptomoneda), concat(min(fecha), ',', max(fecha))
func (r RepositorioCotizaciones) Resumen(parametros types.Filter) (string, string, error) {
	q := QueryBaseCotizaciones(parametros)
	q.Select = "SELECT JSON_ARRAYAGG(simbolo), JSON_OBJECT('fecha_min', min(fecha), 'fecha_max', max(fecha), 'count', count(*))"

	var monedas string
	var fechas string

	if err := r.db.QueryRow(q.toString()).Scan(&monedas, &fechas); err != nil {
		return "", "", err
	}

	return monedas, fechas, nil
}

func (r RepositorioCotizaciones) monedaPorId(id int) (domain.Criptomoneda, error) {
	query := "SELECT id, nombre, simbolo FROM go.criptomoneda WHERE id = ?"

	var moneda domain.Criptomoneda

	err := r.db.QueryRow(query, id).Scan(&moneda.ID, &moneda.Nombre, &moneda.Simbolo)
	if err != nil {
		return domain.Criptomoneda{}, err
	}

	return moneda, nil
}

func (r RepositorioCotizaciones) EsCotizacionManual(id int) (bool, error) {
	query := "SELECT COUNT(*) FROM go.cotizacion WHERE id = ? AND api = 'manual'"
	var rowCount int
	if err := r.db.QueryRow(query, id).Scan(&rowCount); err != nil {
		return false, err
	}

	if rowCount == 0 {
		return false, nil
	}

	return true, nil
}
