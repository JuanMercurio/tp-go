package servicios

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestCriptomonedaService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// repo := mock_ports.NewMockRepositorioMonedas(ctrl)
	// cot := mock_ports.NewMockCotizador(ctrl)

	// assert.NoError(t, err)
	// assert.Equal(t, expectedResult, result)

	// repo.AssertCalled(t, "GetAll")
}
