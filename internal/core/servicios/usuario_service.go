package servicios

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/ports"
)

type UsuarioServicio struct {
	ru ports.RepositorioUsuarios
	rm ports.RepositorioMonedas
}

var (
	ErrDuplicateMail     = errors.New("duplicate email")
	ErrDuplicateUserName = errors.New("duplicate username")
)

func CrearServicioUsuario(ru ports.RepositorioUsuarios, rm ports.RepositorioMonedas) *UsuarioServicio {
	return &UsuarioServicio{
		ru: ru,
		rm: rm,
	}
}

func (s UsuarioServicio) AltaUsuario(params ports.AltaUsuarioParams) (int, error) {

	usuario, err := CrearUsuarioDeParams(params)
	if err != nil {
		return 0, err
	}

	id, err := s.ru.AltaUsuario(usuario)
	if err != nil {
		if strings.Contains(err.Error(), "email") {
			return 0, ErrDuplicateMail
		}

		if strings.Contains(err.Error(), "username") {
			return 0, ErrDuplicateUserName
		}
		return 0, err
	}
	return id, nil
}

func (s UsuarioServicio) BajaUsuario(id int) error {
	return s.ru.BajaUsuario(id)
}

func (s UsuarioServicio) BuscarTodos() ([]domain.Usuario, error) {
	return s.ru.BuscarTodos()
}

func CrearUsuarioDeParams(p ports.AltaUsuarioParams) (domain.Usuario, error) {
	doc, err := domain.CreateDocumento(p.Documento, p.TipoDocumento)
	if err != nil {
		return domain.Usuario{}, fmt.Errorf("error creando el usuario: %w", err)
	}

	return domain.Usuario{
		Username:          p.Username,
		Nombre:            p.Nombre,
		Apellido:          p.Apellido,
		Email:             p.Email,
		FechaDeNacimiento: p.FechaDeNacimiento,
		Documento:         doc,
	}, nil
}

func (s UsuarioServicio) PatchUsuario(id int, patchs []ports.Patch) (domain.Usuario, error) {

	usuario, err := s.ru.UsuarioPorId(id)
	if err != nil {
		return domain.Usuario{}, err
	}

	mapPatchs := make(map[string]string)
	for _, patch := range patchs {
		switch patch.Path {
		case "monedas":
			if err := s.ActualizarMonedasUsuario(usuario, patch.Op, patch.NuevoValor.(string)); err != nil {
				return domain.Usuario{}, err
			}

		default:
			mapPatchs[patch.Path] = patch.NuevoValor.(string)
		}
	}

	if err := s.ru.ActualizarUsuarioConMap(id, mapPatchs); err != nil {
		return domain.Usuario{}, err
	}

	return s.ru.UsuarioPorId(id)
}

func (s UsuarioServicio) ActualizarMonedasUsuario(usuario domain.Usuario, op string, valor string) error {

	noWhitespace := strings.ReplaceAll(valor, " ", "")
	simbolos := strings.Split(noWhitespace, ",")

	idsSimbolos, err := s.rm.IdsDeSimbolos(simbolos)
	if err != nil {
		return err
	}

	switch op {
	case "cambiar":
		if err := s.ru.ReemplazarMonedas(usuario, idsSimbolos); err != nil {
			return err
		}

	case "agregar":

		if err := s.ru.AgregarMonedasAUsuario(usuario, idsSimbolos); err != nil {
			return err
		}

	case "quitar":

		if err := s.ru.EliminarMonedasUsuario(usuario.Id, idsSimbolos); err != nil {
			return err
		}

	}
	return nil
}

func (s UsuarioServicio) UsuarioValido(id int) error {
	_, err := s.ru.UsuarioPorId(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("usuario no existe")
		}
		return err
	}

	return nil
}

func (s UsuarioServicio) IdDeUsername(username string) (int, error) {
	return s.ru.IdDeUsername(username)
}
