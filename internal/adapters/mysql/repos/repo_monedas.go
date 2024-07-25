package repos

import (
	"database/sql"
	"errors"
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

	id, err := r.darDeAltaEntidad(query, cotizacion.Moneda.ID, cotizacion.Time, cotizacion.Valor, cotizacion.Api)
	if err != nil {
		return 0, err
	}

	return id, nil
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

// la version facil del select seria --> select group_concat(distinct id_criptomoneda), concat(min(fecha), ',', max(fecha))
func (r RepositorioMoneda) Resumen(parametros ports.Filter) (string, string, error) {
	q := QueryBaseCotizaciones(parametros)
	q.Select = "SELECT JSON_ARRAYAGG(simbolo), JSON_OBJECT('fecha_min', min(fecha), 'fecha_max', max(fecha), 'count', count(*))"

	var monedas string
	var fechas string

	if err := r.db.QueryRow(q.toString()).Scan(&monedas, &fechas); err != nil {
		return "", "", err
	}

	return monedas, fechas, nil
}

func (r RepositorioMoneda) SimbolosValido(simbolos []string) error {

	for _, simbolo := range simbolos {
		if err := r.SimboloValido(simbolo); err != nil {
			return err
		}
	}
	return nil
}

func (r RepositorioMoneda) SimboloValido(simbolo string) error {

	fmt.Printf("El simbolo que recive es: %s\n", simbolo)

	query := "SELECT COUNT(*) FROM go.criptomoneda WHERE simbolo IN (?)"
	var cantidad int

	err := r.db.QueryRow(query, simbolo).Scan(&cantidad)
	if err != nil {
		return err
	}

	if cantidad == 0 {
		return fmt.Errorf("la moneda %s no existe", simbolo)
	}

	return nil
}

func (r RepositorioMoneda) IdDeSimbolo(simbolo string) (int, error) {
	query := ` SELECT id FROM criptomoneda WHERE simbolo = ?  `
	var id int
	if err := r.db.QueryRow(query, simbolo).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// podria ser mas perfomante
func (r RepositorioMoneda) IdsDeSimbolos(simbolos []string) ([]int, error) {
	ids := make([]int, 0)
	for _, simbolo := range simbolos {
		id, err := r.IdDeSimbolo(simbolo)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
func (r RepositorioMoneda) BajaCotizacionManual(id int) error {
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

func (r RepositorioMoneda) CotizacionPorId(id int) (domain.Cotizacion, error) {

	query := "SELECT id_criptomoneda, fecha, valor, api FROM go.cotizacion WHERE id = ?"
	var cotizacion domain.Cotizacion
	var idMoneda int
	if err := r.db.QueryRow(query, id).Scan(&idMoneda, &cotizacion.Time, &cotizacion.Valor, &cotizacion.Api); err != nil {
		return domain.Cotizacion{}, err
	}
	moneda, err := r.BuscarPorId(idMoneda)
	if err != nil {
		return domain.Cotizacion{}, err
	}

	cotizacion.Moneda = moneda
	cotizacion.ID = id

	return cotizacion, nil
}

func (r RepositorioMoneda) ActualizarCotizacionMap2(usuario int, idCotizacion int, cambios map[string]any) error {
	return nil
}

func (r RepositorioMoneda) ActualizarCotizacionMap(usuario int, idCotizacion int, cambios map[string]any) error {
	// r.ActualizarCotizacionMap2(usuario, idCotizacion, cambios)
	// return nil

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := "UPDATE cotizacion SET "

	var sets []string
	for key, value := range cambios {
		//no se puede cambiar el id
		if key == "id" {
			continue
		}
		sets = append(sets, key+" = "+"'"+value.(string)+"'")
	}

	query += strings.Join(sets, ",")
	query += fmt.Sprintf(" WHERE id = %d AND api = 'manual'", idCotizacion)

	result, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	afectadas, err := result.RowsAffected()
	if err != nil {
		return err
	}

	//mejores errores
	if afectadas == 0 {
		query := "SELECT COUNT(*) FROM go.cotizacion WHERE id = ?"
		var count int
		err := r.db.QueryRow(query, idCotizacion).Scan(&count)
		if err != nil {
			return err
		}

		if count == 0 {
			return errors.New("la cotizacion no existe")
		}

		// falta el caso en el que es igual (por ahora devuelve esto)
		return errors.New("la cotizacion no es manual")
	}

	for key, value := range cambios {

		if err := r.Auditar(usuario, idCotizacion, "actualizar", key, value); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	return nil
}

func (r RepositorioMoneda) Auditar(usuario int, cotizacion int, accion string, columnaAfectada string, nuevoValor any) error {
	query := `
	INSERT INTO auditoria (id_usuario, id_cotizacion, accion, columna_afectada, nuevo_valor, fecha) VALUES (?,?,?,?,?,?)	
	`
	res, err := r.db.Exec(query, usuario, cotizacion, accion, columnaAfectada, nuevoValor, time.Now())
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
