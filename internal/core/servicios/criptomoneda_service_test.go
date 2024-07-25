package servicios

// import (
// 	"errors"
// 	"testing"

// 	"github.com/juanmercurio/tp-go/internal/core/domain"
// 	mock_ports "github.com/juanmercurio/tp-go/internal/ports/mock"
// 	"github.com/stretchr/testify/assert"
// 	"go.uber.org/mock/gomock"
// )

// func TestCriptomonedaService_Cotizador(t *testing.T) {
// 	//todo juntar tests por tipo de assertion

// 	casosDePrueba := []struct {
// 		Nombre        string
// 		Api           string
// 		MockCotizador func(string, *mock_ports.MockCotizadorMockRecorder)
// 		Assertions    func(*testing.T, int, error)
// 	}{
// 		{
// 			Nombre: "Error en algun cotizador",
// 			Api:    "api",
// 			MockCotizador: func(moneda string, recorder *mock_ports.MockCotizadorMockRecorder) {
// 				recorder.ExisteMoneda(moneda).AnyTimes().Return(false, errors.New("error"))
// 			},
// 			Assertions: func(t *testing.T, id int, err error) {
// 				assert.NotNil(t, err)
// 			},
// 		},
// 		{
// 			Nombre: "No existe la moneda en algun cotizador",
// 			Api:    "api",
// 			MockCotizador: func(moneda string, recorder *mock_ports.MockCotizadorMockRecorder) {
// 				recorder.ExisteMoneda(moneda).AnyTimes().Return(false, nil)
// 			},
// 			Assertions: func(t *testing.T, id int, err error) {
// 				assert.NotNil(t, err)
// 				assert.Zero(t, id)
// 			},
// 		},
// 		{
// 			Nombre: "Existe la moneda en todos lo cotizadores",
// 			Api:    "api",
// 			MockCotizador: func(moneda string, recorder *mock_ports.MockCotizadorMockRecorder) {
// 				recorder.ExisteMoneda(moneda).AnyTimes().Return(true, nil)
// 			},
// 			Assertions: func(t *testing.T, id int, err error) {
// 				assert.Nil(t, err)
// 				assert.EqualValues(t, id, 1)
// 			},
// 		},
// 	}

// 	for _, caso := range casosDePrueba {

// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		api := "api-cualquiera"
// 		moneda := "nombreMoneda"
// 		simbolo := "SIMBOLO"

// 		cot := mock_ports.NewMockCotizador(ctrl)
// 		cot.EXPECT().GetNombre().AnyTimes().Return(api)

// 		repo := mock_ports.NewMockRepositorioMonedas(ctrl)
// 		repo.EXPECT().AltaMoneda(domain.Criptomoneda{ID: 0, Nombre: "nombreMoneda", Simbolo: "SIMBOLO"}).AnyTimes().Return(1, nil)

// 		caso.MockCotizador(simbolo, cot.EXPECT())

// 		servicio := CrearServicioMoneda(repo, cot)
// 		id, err := servicio.AltaMoneda(moneda, simbolo)
// 		caso.Assertions(t, id, err)

// 	}
// }
