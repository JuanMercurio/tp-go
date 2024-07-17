package ports

type CotizacionOutputDTO struct {
	NombreMoneda string
	Simbolo      string
	Fecha        string
	Valor        float64
	Api          string
}

type MonedaDTOOutput struct {
	Id           int
	NombreMoneda string
	Simbolo      string
}
