package types

type CotizacionPut struct {
	Fecha   string  `json:"fecha"`
	Simbolo string  `json:"simbolo"`
	Valor   float64 `json:"valor"`
}
