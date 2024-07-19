package servicios

import (
	"fmt"
	"strings"

	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/ports"
)

type UsuarioServicio struct {
	repo ports.RepositorioUsuarios
}

func CrearServicioUsuario(repo ports.RepositorioUsuarios) *UsuarioServicio {
	return &UsuarioServicio{
		repo: repo,
	}
}

func (s UsuarioServicio) AltaUsuario(nombre string) (int, error) {
	id, err := s.repo.AltaUsuario(domain.CrearUsuario(nombre))
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return 0, fmt.Errorf("el usuario ya existe")
		}

		return 0, err
	}
	return id, nil
}

func (s UsuarioServicio) BajaUsuario(id int) error {
	return s.repo.BajaUsuario(id)
}

func (s UsuarioServicio) BuscarTodos() ([]domain.Usuario, error) {
	return s.repo.BuscarTodos()
}
