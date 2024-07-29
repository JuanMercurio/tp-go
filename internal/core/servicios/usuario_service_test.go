package servicios

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/juanmercurio/tp-go/internal/core/domain"
	mock_ports "github.com/juanmercurio/tp-go/internal/ports/mock"
	"github.com/juanmercurio/tp-go/internal/ports/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAltaUsuario_Success(t *testing.T) {

	params := types.AltaUsuarioParams{
		Username:          "username",
		Nombre:            "nombre",
		Apellido:          "apellido",
		FechaDeNacimiento: time.Now(),
		Documento:         "documento",
		TipoDocumento:     "DNI",
		Email:             "email",
	}

	ctrl := gomock.NewController(t)
	ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
	ru.EXPECT().AltaUsuario(gomock.Any()).AnyTimes().Return(1, nil)

	id, err := CrearServicioUsuario(ru, nil).AltaUsuario(params)
	assert.Nil(t, err)
	assert.Equal(t, id, 1)
}

func TestAltaUsuario_Fail(t *testing.T) {

	paramsValidos := types.AltaUsuarioParams{
		Username:          "username",
		Nombre:            "nombre",
		Apellido:          "apellido",
		FechaDeNacimiento: time.Now(),
		Documento:         "documento",
		TipoDocumento:     "DNI",
		Email:             "email",
	}

	ErrOtro := errors.New("otro error")

	testCases := []struct {
		name        string
		params      types.AltaUsuarioParams
		useCase     *UsuarioServicio
		expectedErr error
	}{
		{
			name: "Error de tipo de dni incorrecto",
			params: types.AltaUsuarioParams{
				Username:          "username",
				Nombre:            "nombre",
				Apellido:          "apellido",
				FechaDeNacimiento: time.Now(),
				Documento:         "documento",
				TipoDocumento:     "mal tipo dni",
				Email:             "email",
			},
			useCase: func() *UsuarioServicio {
				ctrl := gomock.NewController(t)
				ru := mock_ports.NewMockRepositorioUsuarios(ctrl)

				return CrearServicioUsuario(ru, nil)
			}(),
			expectedErr: domain.ErrTipoDocumentoInexistente,
		},
		{
			name:   "Error cuando ya existe el mail",
			params: paramsValidos,
			useCase: func() *UsuarioServicio {
				ctrl := gomock.NewController(t)
				ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
				ru.EXPECT().AltaUsuario(gomock.Any()).AnyTimes().Return(0, fmt.Errorf("email duplicado"))

				return CrearServicioUsuario(ru, nil)
			}(),
			expectedErr: ErrDuplicateMail,
		},
		{
			name:   "Error el usuario ya existe",
			params: paramsValidos,
			useCase: func() *UsuarioServicio {
				ctrl := gomock.NewController(t)
				ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
				ru.EXPECT().AltaUsuario(gomock.Any()).Return(0, fmt.Errorf("username duplicado"))

				return CrearServicioUsuario(ru, nil)
			}(),
			expectedErr: ErrDuplicateUserName,
		},
		{
			name:   "Error inesperado de alta de usuario",
			params: paramsValidos,
			useCase: func() *UsuarioServicio {
				ctrl := gomock.NewController(t)
				ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
				ru.EXPECT().AltaUsuario(gomock.Any()).Return(0, ErrOtro)

				return CrearServicioUsuario(ru, nil)
			}(),
			expectedErr: ErrOtro,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := testCase.useCase.AltaUsuario(testCase.params)

			assertions := assert.New(t)
			assertions.True(errors.Is(err, testCase.expectedErr))
		})
	}
}

func TestPatchUsuario_Fail(t *testing.T) {

	patches := []types.Patch{
		{
			Op:         "cambiar",
			Path:       "monedas",
			NuevoValor: "BTC, ASDF",
		},
		{
			Op:         "quitar",
			Path:       "monedas",
			NuevoValor: "",
		},
		{
			Op:         "agregar",
			Path:       "monedas",
			NuevoValor: "BTC, ASDF",
		},
	}

	moneda := domain.Criptomoneda{
		Nombre:  "nombre",
		Simbolo: "simbolo",
		ID:      1,
	}

	doc := domain.Documento{
		Tipo:   domain.DNI,
		Numero: "numero",
	}

	usuario := domain.Usuario{
		Id:                1,
		Username:          "user",
		Nombre:            "nombre",
		Apellido:          "apellido",
		FechaDeNacimiento: time.Now(),
		Documento:         doc,
		Email:             "email@email.com",
		MonedasInteres:    []domain.Criptomoneda{moneda},
	}

	testCases := []struct {
		name        string
		useCase     *UsuarioServicio
		expectedErr string
	}{
		{
			name: "El id del usuario no existe",
			useCase: func() *UsuarioServicio {
				ctrl := gomock.NewController(t)
				ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
				ru.EXPECT().UsuarioPorId(1).Return(domain.Usuario{}, fmt.Errorf("no existe usuario"))
				return CrearServicioUsuario(ru, nil)
			}(),
			expectedErr: "no existe usuario",
		},
		{
			name: "Se quiere cambiar a una moneda que no esta registrada",
			useCase: func() *UsuarioServicio {
				ctrl := gomock.NewController(t)
				ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
				rm := mock_ports.NewMockRepositorioMonedas(ctrl)

				ru.EXPECT().UsuarioPorId(1).Return(usuario, nil)
				rm.EXPECT().IdsDeSimbolos(gomock.Any()).Return(nil, fmt.Errorf("no existe moneda"))

				return CrearServicioUsuario(ru, rm)
			}(),
			expectedErr: "no existe moneda",
		},
		{
			name: "Error reemplazando monedas",
			useCase: func() *UsuarioServicio {
				monedas := []int{1}
				ctrl := gomock.NewController(t)

				ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
				rm := mock_ports.NewMockRepositorioMonedas(ctrl)

				ru.EXPECT().UsuarioPorId(1).Return(usuario, nil)
				rm.EXPECT().IdsDeSimbolos(gomock.Any()).Return(monedas, nil)
				ru.EXPECT().ReemplazarMonedas(gomock.Any(), gomock.Any()).Return(errors.New("error reemplazando monedas"))

				return CrearServicioUsuario(ru, rm)
			}(),
			expectedErr: "error reemplazando monedas",
		},
		{
			name: "Error borrando monedas",
			useCase: func() *UsuarioServicio {
				monedas := []int{1}
				ctrl := gomock.NewController(t)

				ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
				rm := mock_ports.NewMockRepositorioMonedas(ctrl)

				ru.EXPECT().UsuarioPorId(1).Return(usuario, nil)
				rm.EXPECT().IdsDeSimbolos(gomock.Any()).AnyTimes().Return(monedas, nil)
				ru.EXPECT().ReemplazarMonedas(gomock.Any(), gomock.Any()).Return(nil)
				ru.EXPECT().EliminarMonedasUsuario(gomock.Any(), gomock.Any()).Return(errors.New("error eliminando monedas"))

				return CrearServicioUsuario(ru, rm)
			}(),
			expectedErr: "error eliminando monedas",
		},
		{
			name: "Error agregando monedas",
			useCase: func() *UsuarioServicio {
				monedas := []int{1}
				ctrl := gomock.NewController(t)

				ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
				rm := mock_ports.NewMockRepositorioMonedas(ctrl)

				ru.EXPECT().UsuarioPorId(1).Return(usuario, nil)
				rm.EXPECT().IdsDeSimbolos(gomock.Any()).AnyTimes().Return(monedas, nil)
				ru.EXPECT().ReemplazarMonedas(gomock.Any(), gomock.Any()).Return(nil)
				ru.EXPECT().EliminarMonedasUsuario(gomock.Any(), gomock.Any()).Return(nil)
				ru.EXPECT().AgregarMonedasAUsuario(gomock.Any(), gomock.Any()).Return(errors.New("error agregando monedas"))

				return CrearServicioUsuario(ru, rm)
			}(),
			expectedErr: "error agregando monedas",
		},
		{
			name: "Error actualizando campos con mapa",
			useCase: func() *UsuarioServicio {
				monedas := []int{1}
				ctrl := gomock.NewController(t)

				ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
				rm := mock_ports.NewMockRepositorioMonedas(ctrl)

				ru.EXPECT().UsuarioPorId(1).Return(usuario, nil)
				rm.EXPECT().IdsDeSimbolos(gomock.Any()).AnyTimes().Return(monedas, nil)
				ru.EXPECT().ReemplazarMonedas(gomock.Any(), gomock.Any()).Return(nil)
				ru.EXPECT().EliminarMonedasUsuario(gomock.Any(), gomock.Any()).Return(nil)
				ru.EXPECT().AgregarMonedasAUsuario(gomock.Any(), gomock.Any()).Return(nil)
				ru.EXPECT().ActualizarUsuarioConMap(gomock.Any(), gomock.Any()).Return(errors.New("error de actualizacion"))

				return CrearServicioUsuario(ru, rm)
			}(),
			expectedErr: "error de actualizacion",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := testCase.useCase.PatchUsuario(1, patches)

			assertions := assert.New(t)
			assertions.True(strings.Contains(err.Error(), testCase.expectedErr))
		})
	}

}

func TestPatchUsuario_Success(t *testing.T) {

	patches := []types.Patch{
		{
			Op:         "cambiar",
			Path:       "monedas",
			NuevoValor: "BTC, ASDF",
		},
		{
			Op:         "quitar",
			Path:       "monedas",
			NuevoValor: "",
		},
		{
			Op:         "agregar",
			Path:       "monedas",
			NuevoValor: "BTC, ASDF",
		},
	}

	moneda := domain.Criptomoneda{
		Nombre:  "nombre",
		Simbolo: "simbolo",
		ID:      1,
	}

	doc := domain.Documento{
		Tipo:   domain.DNI,
		Numero: "numero",
	}

	usuario := domain.Usuario{
		Id:                1,
		Username:          "user",
		Nombre:            "nombre",
		Apellido:          "apellido",
		FechaDeNacimiento: time.Now(),
		Documento:         doc,
		Email:             "email@email.com",
		MonedasInteres:    []domain.Criptomoneda{moneda},
	}

	monedas := []int{1}
	ctrl := gomock.NewController(t)

	ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
	rm := mock_ports.NewMockRepositorioMonedas(ctrl)

	ru.EXPECT().UsuarioPorId(1).Return(usuario, nil)
	rm.EXPECT().IdsDeSimbolos(gomock.Any()).AnyTimes().Return(monedas, nil)
	ru.EXPECT().ReemplazarMonedas(gomock.Any(), gomock.Any()).Return(nil)
	ru.EXPECT().EliminarMonedasUsuario(gomock.Any(), gomock.Any()).Return(nil)
	ru.EXPECT().AgregarMonedasAUsuario(gomock.Any(), gomock.Any()).Return(nil)
	ru.EXPECT().ActualizarUsuarioConMap(gomock.Any(), gomock.Any()).Return(nil)

	_, err := CrearServicioUsuario(ru, rm).PatchUsuario(1, patches)
	assert.Nil(t, err)
}

func TestCrearUsuario_Fail(t *testing.T) {

	paramsTipoDocInvalido := types.AltaUsuarioParams{
		Username:          "username",
		Nombre:            "nombre",
		Apellido:          "apellido",
		FechaDeNacimiento: time.Now(),
		Documento:         "documento",
		TipoDocumento:     "ESTE TIPO DE DNI NO EXISTE",
		Email:             "email",
	}

	_, err := CrearUsuarioDeParams(paramsTipoDocInvalido)

	assert.True(t, errors.Is(err, domain.ErrTipoDocumentoInexistente))
}

func TestCrearUsuario_Success(t *testing.T) {

	paramsTipoDocInvalido := types.AltaUsuarioParams{
		Username:          "username",
		Nombre:            "nombre",
		Apellido:          "apellido",
		FechaDeNacimiento: time.Now(),
		Documento:         "documento",
		Email:             "email",
	}

	tiposDoc := []string{"DNI", "CEDULA", "PASAPORTE", "dni", "cedula", "pasaporte"}

	for _, tipo := range tiposDoc {
		t.Run("Crear con usuario con tipo de documento: "+tipo, func(t *testing.T) {
			paramsTipoDocInvalido.TipoDocumento = tipo
			_, err := CrearUsuarioDeParams(paramsTipoDocInvalido)

			assert.Nil(t, err)
		})
	}
}

func TestBajaUsuario_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
	ru.EXPECT().BajaUsuario(1).Return(nil)
	err := CrearServicioUsuario(ru, nil).BajaUsuario(1)
	assert.Nil(t, err)
}
func TestBajaUsuario_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
	ru.EXPECT().BajaUsuario(1).Return(errors.New("error en el repo haciendo la baja"))
	err := CrearServicioUsuario(ru, nil).BajaUsuario(1)
	assert.NotNil(t, err)

}
func TestBuscarTodos_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
	ru.EXPECT().BuscarTodos().Return([]domain.Usuario{}, nil)
	_, err := CrearServicioUsuario(ru, nil).BuscarTodos()
	assert.Nil(t, err)
}

func TestBuscarTodos_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	ru := mock_ports.NewMockRepositorioUsuarios(ctrl)
	ru.EXPECT().BuscarTodos().Return([]domain.Usuario{}, errors.New("error en el repo buscando todos"))
	_, err := CrearServicioUsuario(ru, nil).BuscarTodos()
	assert.NotNil(t, err)
}

func TestCrearServicioUsuario(t *testing.T) {
	ru := mock_ports.NewMockRepositorioUsuarios(gomock.NewController(t))
	rm := mock_ports.NewMockRepositorioMonedas(gomock.NewController(t))
	s := CrearServicioUsuario(ru, rm)
	assert.NotNil(t, s)
}
