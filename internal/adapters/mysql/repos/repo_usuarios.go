package repos

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

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

func (r RepositorioUsuario) AltaUsuario(u domain.Usuario) (int, error) {
	query := "INSERT INTO usuario  (username, nombre, apellido, email, doc, tipo_doc, fecha_nacimiento, activo) VALUES (?,?,?,?,?,?, ?, true)"

	tipodoc, documento := u.Documento.Split()

	result, err := r.db.Exec(query, u.Username, u.Nombre, u.Apellido, u.Email, documento, tipodoc, u.FechaDeNacimiento)
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

func (r RepositorioUsuario) MonedasDeUsuario(id int) ([]domain.Criptomoneda, error) {

	var builder QueryBuilder
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

func (r RepositorioUsuario) BuscarTodos() ([]domain.Usuario, error) {
	query := "SELECT id, username,  nombre, apellido, email, doc, tipo_doc, fecha_nacimiento FROM usuario WHERE activo is not  false"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error eejcutando la query: %w", err)
	}

	var u domain.Usuario
	var usuarios []domain.Usuario
	var docString, tipoString, fechaString string

	for rows.Next() {
		if err := rows.Scan(&u.Id, &u.Username, &u.Nombre, &u.Apellido, &u.Email, &docString, &tipoString, &fechaString); err != nil {
			return nil, fmt.Errorf("error escanenado las rows: %w", err)
		}

		fecha, err := time.Parse("2006-01-02 15:04:05", fechaString)
		if err != nil {
			return nil, err
		}
		u.FechaDeNacimiento = fecha

		u.Documento, err = domain.CreateDocumento(docString, tipoString)
		if err != nil {
			return nil, err
		}

		u.MonedasInteres, err = r.MonedasDeUsuario(u.Id)
		if err != nil {
			return nil, err
		}
		fmt.Println(u)

		usuarios = append(usuarios, u)
	}
	return usuarios, nil
}

func (r RepositorioUsuario) IsDuplicateString(tabla, columna string, valor string) (bool, error) {

	query := "SELECT COUNT(*) FROM ? WHERE ? = ?"
	var count int

	if err := r.db.QueryRow(query, tabla, columna, valor).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r RepositorioUsuario) UsuarioPorId(id int) (domain.Usuario, error) {
	query := "SELECT username, apellido, doc, tipo_doc, email, fecha_nacimiento FROM usuario WHERE id = ?"

	var u domain.Usuario
	var fechaString string
	var doc, tipo_doc string

	if err := r.db.QueryRow(query, id).Scan(&u.Nombre, &u.Apellido, &doc, &tipo_doc, &u.Email, &fechaString); err != nil {
		return domain.Usuario{}, err
	}

	var err error
	u.MonedasInteres, err = r.MonedasDeUsuario(id)
	if err != nil {
		return domain.Usuario{}, err
	}

	u.FechaDeNacimiento, err = time.Parse("2006-01-02 15:04:05", fechaString)
	if err != nil {
		return domain.Usuario{}, err
	}

	u.Id = id
	u.Documento, err = domain.CreateDocumento(doc, tipo_doc)
	if err != nil {
		return domain.Usuario{}, err
	}

	return u, nil
}

func (r RepositorioUsuario) ActualizarUsuarioConMap(id int, cambios map[string]string) error {
	if len(cambios) == 0 {
		return nil
	}

	query := "UPDATE usuario SET "
	var sets []string
	for key, value := range cambios {

		//no se puede cambiar el id
		if key == "id" {
			continue
		}
		sets = append(sets, key+" = "+"'"+value+"'")
	}

	query += strings.Join(sets, ",")
	query += fmt.Sprintf(" WHERE id = %d", id)

	fmt.Println(query)
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r RepositorioUsuario) ReemplazarMonedas(usuario domain.Usuario, idMonedas []int) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if err := r.EliminarTodasMonedasUsuario(usuario.Id); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.AgregarMonedasAUsuario(usuario, idMonedas); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r RepositorioUsuario) EliminarTodasMonedasUsuario(id int) error {
	query := fmt.Sprintf("DELETE FROM usuario_criptomoneda WHERE id_usuario = %d", id)

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (r RepositorioUsuario) AgregarMonedasAUsuario(usuario domain.Usuario, idMonedas []int) error {
	if len(idMonedas) == 0 {
		return nil
	}

	query := "INSERT INTO usuario_criptomoneda (id_usuario, id_criptomoneda) VALUES "

	for _, id := range idMonedas {
		query += fmt.Sprintf("(%d, %d),", usuario.Id, id)
	}
	query = strings.TrimRight(query, ",")

	if _, err := r.db.Exec(query); err != nil {
		return err
	}

	return nil
}

func (r RepositorioUsuario) EliminarMonedasUsuario(idUsuario int, idsMonedas []int) error {

	ids := ""
	for _, id := range idsMonedas {
		ids += fmt.Sprintf("%d,", id)
	}
	ids = strings.TrimRight(ids, ",")

	query := fmt.Sprintf("DELETE FROM usuario_criptomoneda WHERE id_usuario = ? AND id_criptomoneda IN (%s)", ids)

	if err := r.db.QueryRow(query, idUsuario).Scan(); err != nil {
		return err
	}

	return nil
}
