package types

type CotizacionOutputDTO struct {
	NombreMoneda string  `json:"nombre"`
	Simbolo      string  `json:"simbolo"`
	Fecha        string  `json:"fecha"`
	Valor        float64 `json:"valor"`
	Api          string  `json:"api"`
}

type MonedaOutputDTO struct {
	Id           int    `json:"id"`
	NombreMoneda string `json:"nombre"`
	Simbolo      string `json:"simbolo"`
}
