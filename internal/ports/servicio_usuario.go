package ports

import (
	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/ports/types"
)

type ServicioUsuarios interface {
	AltaUsuario(types.AltaUsuarioParams) (int, error)
	BajaUsuario(id int) error
	BuscarTodos() ([]domain.Usuario, error)
	PatchUsuario(id int, patch []types.Patch) (domain.Usuario, error)
	UsuarioValido(id int) error
	IdDeUsername(username string) (int, error)
}
