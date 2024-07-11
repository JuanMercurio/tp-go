package domain

type Criptomoneda struct {
	ID     int    `json:"id"`
	Nombre string `json:"string"`
}

func New(nombre string) Criptomoneda {
	return Criptomoneda{
		Nombre: nombre,
	}
}
