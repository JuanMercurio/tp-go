package ports

import "github.com/juanmercurio/tp-go/internal/core/domain"

type RepositorioUsuarios interface {
	AltaUsuario(domain.Usuario) (int, error)
	BajaUsuario(int) error
	BuscarTodos() ([]domain.Usuario, error)
}
