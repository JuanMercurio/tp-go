package domain

import "strings"

type Criptomoneda struct {
	ID      int    `json:"id"`
	Nombre  string `json:"string"`
	Simbolo string `json:"simbolo"`
}

func CrearMoneda(nombre, simbolo string) Criptomoneda {
	return Criptomoneda{
		Nombre:  nombre,
		Simbolo: strings.ToUpper(simbolo),
	}
}
