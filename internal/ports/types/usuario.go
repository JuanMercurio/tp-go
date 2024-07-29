package types

import "time"

type AltaUsuarioParams struct {
	Username          string
	Nombre            string
	Apellido          string
	FechaDeNacimiento time.Time
	Documento         string
	TipoDocumento     string
	Email             string
}
