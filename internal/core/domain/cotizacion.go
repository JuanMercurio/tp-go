package domain

import "time"

type Cotizacion struct {
	ID     int          `json:"id"`
	Moneda Criptomoneda `json:"criptomoneda"`
	Valor  float64      `json:"valor"`
	Time   time.Time    `json:"fecha"`
	Api    string       `json:"api"`
}
