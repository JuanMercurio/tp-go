// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juanmercurio/tp-go/internal/ports (interfaces: RepositorioMonedas)
//
// Generated by this command:
//
//	mockgen --destination=./mock/repo_monedas.go github.com/juanmercurio/tp-go/internal/ports RepositorioMonedas
//

// Package mock_ports is a generated GoMock package.
package mock_ports

import (
	reflect "reflect"

	domain "github.com/juanmercurio/tp-go/internal/core/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockRepositorioMonedas is a mock of RepositorioMonedas interface.
type MockRepositorioMonedas struct {
	ctrl     *gomock.Controller
	recorder *MockRepositorioMonedasMockRecorder
}

// MockRepositorioMonedasMockRecorder is the mock recorder for MockRepositorioMonedas.
type MockRepositorioMonedasMockRecorder struct {
	mock *MockRepositorioMonedas
}

// NewMockRepositorioMonedas creates a new mock instance.
func NewMockRepositorioMonedas(ctrl *gomock.Controller) *MockRepositorioMonedas {
	mock := &MockRepositorioMonedas{ctrl: ctrl}
	mock.recorder = &MockRepositorioMonedasMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositorioMonedas) EXPECT() *MockRepositorioMonedasMockRecorder {
	return m.recorder
}

// AltaMoneda mocks base method.
func (m *MockRepositorioMonedas) AltaMoneda(arg0 domain.Criptomoneda) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AltaMoneda", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AltaMoneda indicates an expected call of AltaMoneda.
func (mr *MockRepositorioMonedasMockRecorder) AltaMoneda(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AltaMoneda", reflect.TypeOf((*MockRepositorioMonedas)(nil).AltaMoneda), arg0)
}

// BuscarPorId mocks base method.
func (m *MockRepositorioMonedas) BuscarPorId(arg0 int) (domain.Criptomoneda, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuscarPorId", arg0)
	ret0, _ := ret[0].(domain.Criptomoneda)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuscarPorId indicates an expected call of BuscarPorId.
func (mr *MockRepositorioMonedasMockRecorder) BuscarPorId(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuscarPorId", reflect.TypeOf((*MockRepositorioMonedas)(nil).BuscarPorId), arg0)
}

// BuscarTodos mocks base method.
func (m *MockRepositorioMonedas) BuscarTodos() ([]domain.Criptomoneda, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuscarTodos")
	ret0, _ := ret[0].([]domain.Criptomoneda)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuscarTodos indicates an expected call of BuscarTodos.
func (mr *MockRepositorioMonedasMockRecorder) BuscarTodos() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuscarTodos", reflect.TypeOf((*MockRepositorioMonedas)(nil).BuscarTodos))
}

// IdDeSimbolo mocks base method.
func (m *MockRepositorioMonedas) IdDeSimbolo(arg0 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IdDeSimbolo", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IdDeSimbolo indicates an expected call of IdDeSimbolo.
func (mr *MockRepositorioMonedasMockRecorder) IdDeSimbolo(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IdDeSimbolo", reflect.TypeOf((*MockRepositorioMonedas)(nil).IdDeSimbolo), arg0)
}

// IdsDeSimbolos mocks base method.
func (m *MockRepositorioMonedas) IdsDeSimbolos(arg0 []string) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IdsDeSimbolos", arg0)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IdsDeSimbolos indicates an expected call of IdsDeSimbolos.
func (mr *MockRepositorioMonedasMockRecorder) IdsDeSimbolos(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IdsDeSimbolos", reflect.TypeOf((*MockRepositorioMonedas)(nil).IdsDeSimbolos), arg0)
}

// SimboloValido mocks base method.
func (m *MockRepositorioMonedas) SimboloValido(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SimboloValido", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SimboloValido indicates an expected call of SimboloValido.
func (mr *MockRepositorioMonedasMockRecorder) SimboloValido(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SimboloValido", reflect.TypeOf((*MockRepositorioMonedas)(nil).SimboloValido), arg0)
}

// SimbolosValido mocks base method.
func (m *MockRepositorioMonedas) SimbolosValido(arg0 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SimbolosValido", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SimbolosValido indicates an expected call of SimbolosValido.
func (mr *MockRepositorioMonedasMockRecorder) SimbolosValido(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SimbolosValido", reflect.TypeOf((*MockRepositorioMonedas)(nil).SimbolosValido), arg0)
}
