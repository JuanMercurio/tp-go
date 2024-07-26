package ports

import (
	"github.com/juanmercurio/tp-go/internal/core/domain"
)

type ServicioUsuarios interface {
	AltaUsuario(AltaUsuarioParams) (int, error)
	BajaUsuario(id int) error
	BuscarTodos() ([]domain.Usuario, error)
	PatchUsuario(id int, patch []Patch) (domain.Usuario, error)
	UsuarioValido(id int) error
	IdDeUsername(username string) (int, error)
}
