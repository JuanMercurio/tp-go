package domain

import (
	"errors"
	"strings"
	"time"
)

type Usuario struct {
	Id                int
	Username          string
	Nombre            string
	Apellido          string
	FechaDeNacimiento time.Time
	Documento         Documento
	Email             string
	MonedasInteres    []Criptomoneda
}

type TipoDocumento int

const (
	DNI TipoDocumento = iota
	Cedula
	Pasaporte
)

type Documento struct {
	Tipo   TipoDocumento
	Numero string
}

func (d Documento) Split() (string, string) {
	var tipo string
	switch d.Tipo {
	case DNI:
		tipo = "DNI"
	case Pasaporte:
		tipo = "PASAPORTE"
	case Cedula:
		tipo = "CEDULA"
	}
	return tipo, d.Numero
}

func CreateDocumento(documento, tipo string) (Documento, error) {

	var doc Documento
	doc.Numero = documento

	switch strings.ToUpper(tipo) {
	case "DNI":
		doc.Tipo = DNI
	case "PASAPORTE":
		doc.Tipo = Pasaporte
	case "CEDULA":
		doc.Tipo = Cedula
	default:
		return Documento{}, errors.New("solo se aceptan DNI | PASAPORTE | CEDULA")
	}
	return doc, nil
}
