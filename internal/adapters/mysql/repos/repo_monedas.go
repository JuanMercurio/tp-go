package repos

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/juanmercurio/tp-go/internal/core/domain"
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

// deberia estar afuera de este archivo y ser usada por los demas repos
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

func (r RepositorioMoneda) SimbolosValido(simbolos []string) error {
	if len(simbolos) == 0 || simbolos[0] == "" {
		return nil
	}

	for _, simbolo := range simbolos {
		if err := r.SimboloValido(simbolo); err != nil {
			return err
		}
	}
	return nil
}

func (r RepositorioMoneda) SimboloValido(simbolo string) error {

	query := "SELECT COUNT(*) FROM go.criptomoneda WHERE simbolo IN (?)"
	var cantidad int

	err := r.db.QueryRow(query, simbolo).Scan(&cantidad)
	if err != nil {
		return err
	}

	if cantidad == 0 {
		return fmt.Errorf("la moneda %s no esta registrada", simbolo)
	}

	return nil
}

func (r RepositorioMoneda) IdDeSimbolo(simbolo string) (int, error) {
	query := ` SELECT id FROM criptomoneda WHERE simbolo = ?  `
	var id int
	if err := r.db.QueryRow(query, simbolo).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("la moneda %s no esta registrada: %w", simbolo, err)
		}
		return 0, err
	}
	return id, nil
}

// podria ser mas perfomante
func (r RepositorioMoneda) IdsDeSimbolos(simbolos []string) ([]int, error) {
	if len(simbolos) == 0 || simbolos[0] == "" {
		return nil, nil
	}

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
