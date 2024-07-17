package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrearPaginas(t *testing.T) {
	tamPaginas := 10
	elementos := []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	expected := []Pagina{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11, 12, 13, 14, 15}}
	result := crearPaginas(tamPaginas, elementos)
	assert.Equal(t, expected, result)
}

func TestPaginasNecesarias(t *testing.T) {
	tamPaginas := 10
	filas := 15
	expected := 2
	result := paginasNecesarias(tamPaginas, filas)
	assert.Equal(t, expected, result)
}
func TestPaginasNecesariasConPaginaMasGrande(t *testing.T) {
	tamPaginas := 10
	filas := 2
	expected := 1
	result := paginasNecesarias(tamPaginas, filas)
	assert.Equal(t, expected, result)
}

func TestObtenerSegunoffsetYCantidad(t *testing.T) {
	elementos := []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	offset := 4
	cant := 5
	expected := Pagina{5, 6, 7, 8, 9}
	result := obtenerElementos(elementos, offset, cant)
	assert.Equal(t, expected, result)
}
