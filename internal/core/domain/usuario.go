package domain

type Usuario struct {
	Id     int
	Nombre string
}

func CrearUsuario(nombre string) Usuario {
	return Usuario{
		Nombre: nombre,
	}
}
