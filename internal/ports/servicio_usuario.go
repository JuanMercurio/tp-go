package ports

import "github.com/juanmercurio/tp-go/internal/core/domain"

type ServicioUsuarios interface {
	AltaUsuario(nombre string) (int, error)
	BajaUsuario(id int) error
	BuscarTodos() ([]domain.Usuario, error)
}
