package servicios

import (
	"errors"
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
	id, err := s.ru.AltaUsuario(CrearUsuario(params))
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

func CrearUsuario(p ports.AltaUsuarioParams) domain.Usuario {
	return domain.Usuario{
		Username:          p.Username,
		Nombre:            p.Nombre,
		Apellido:          p.Apellido,
		Email:             p.Email,
		FechaDeNacimiento: p.FechaDeNacimiento,
		Documento:         obtenerDocumento(p.TipoDocumento, p.Documento),
	}
}

// todo redundante hay otra funcion sacar
func obtenerDocumento(t, n string) domain.Documento {
	var doc domain.Documento
	switch t {
	case "DNI":
		doc.Tipo = domain.DNI
	case "CEDULA":
		doc.Tipo = domain.Cedula
	case "PASAPORTE":
		doc.Tipo = domain.Pasaporte
	}

	doc.Numero = n
	return doc
}

func (s UsuarioServicio) PatchUsuario(id int, patchs []ports.Patch) error {

	usuario, err := s.ru.UsuarioPorId(id)
	if err != nil {
		return err
	}

	mapPatchs := make(map[string]string)
	for _, patch := range patchs {
		switch patch.Path {
		case "monedas":
			if err := s.ActualizarMonedasUsuario(usuario, patch.Op, patch.NuevoValor.(string)); err != nil {
				return err
			}

		default:
			mapPatchs[patch.Path] = patch.NuevoValor.(string)
		}
	}

	if err := s.ru.ActualizarUsuarioConMap(id, mapPatchs); err != nil {
		return err
	}

	return nil
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
