//go:generate  /home/feb/go/bin/mockgen  --destination=./mock/repo_usuarios.go github.com/juanmercurio/tp-go/internal/ports RepositorioUsuarios

package ports

import "github.com/juanmercurio/tp-go/internal/core/domain"

type RepositorioUsuarios interface {
	AltaUsuario(usuario domain.Usuario) (int, error)
	BajaUsuario(id int) error
	BuscarTodos() ([]domain.Usuario, error)
	UsuarioPorId(id int) (domain.Usuario, error)
	ActualizarUsuarioConMap(id int, cambios map[string]string) error
	ReemplazarMonedas(usuario domain.Usuario, idsMonedas []int) error
	AgregarMonedasAUsuario(usuario domain.Usuario, idsMonedas []int) error
	EliminarMonedasUsuario(idUsuario int, idsMonedas []int) error
	IdDeUsername(string) (int, error)
}
