package paprika

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type ServicioMockeado struct {
	mock.Mock
}

// GetSimbols is a mock method that returns a list of symbols.
func (m *ServicioMockeado) GetSimbols() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func TestObtenerCotizacion(t *testing.T) {

}
