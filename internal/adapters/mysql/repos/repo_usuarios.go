package repos

import (
	"database/sql"
	"fmt"

	"github.com/juanmercurio/tp-go/internal/core/domain"
)

type RepositorioUsuario struct {
	db *sql.DB
}

func CrearRepositorioUsuario(db *sql.DB) *RepositorioUsuario {
	return &RepositorioUsuario{
		db: db,
	}
}

func (r RepositorioUsuario) AltaUsuario(usuario domain.Usuario) (int, error) {
	query := "INSERT INTO usuario  (nombre, activo ) VALUES ( ?, true)"

	result, err := r.db.Exec(query, usuario.Nombre)
	if err != nil {
		return 0, fmt.Errorf("error al ejecutar el query en la base: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al buscar el id de la moneda insertada: %w", err)
	}

	return int(id), nil
}

// esta es una baja de usuario logica
func (r RepositorioUsuario) BajaUsuario(id int) error {
	query := "UPDATE usuario SET activo = false  where id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	dadosDeBaja, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if dadosDeBaja == 0 {
		return fmt.Errorf("id incorrecto")
	}

	// si se quiere eliminar de la base de datos por ahora no hay un
	// endpoint para hacer eso pero la funcionalidad esta y se deberia
	// usar este metodo
	// r.BajaUsuarioEliminando(id)

	return nil
}

func (r RepositorioUsuario) BajaUsuarioEliminando(id int) error {
	queryRelaciones := "DELETE FROM usuario_criptomoneda WHERE id_usuario = ?"
	queryUsuario := "DELETE FROM usuario WHERE id = ?"

	_, err := r.db.Exec(queryRelaciones, id)
	if err != nil {
		return err
	}

	resultado, err := r.db.Exec(queryUsuario, id)
	if err != nil {
		return err
	}

	eliminados, err := resultado.RowsAffected()
	if err != nil {
		return err
	}

	if eliminados == 0 {
		return fmt.Errorf("id incorrecto")
	}

	return nil
}

func (r RepositorioUsuario) BuscarTodos() ([]domain.Usuario, error) {
	query := "SELECT id, nombre FROM usuario WHERE activo is not  false"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error eejcutando la query: %w", err)
	}

	var usuario domain.Usuario
	var usuarios []domain.Usuario

	for rows.Next() {
		if err := rows.Scan(&usuario.Id, &usuario.Nombre); err != nil {
			return nil, fmt.Errorf("error escanenado las rows: %w", err)
		}
		usuarios = append(usuarios, usuario)
	}
	return usuarios, nil
}
