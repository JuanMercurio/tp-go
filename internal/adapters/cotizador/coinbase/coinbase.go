package coinbase

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/juanmercurio/tp-go/internal/adapters/config"
	"github.com/juanmercurio/tp-go/internal/ports"
)

type Coinbase struct {
	Nombre string
	Url    string
	Token  string
}

func Crear(config *config.APIConfig) *Coinbase {
	return &Coinbase{
		Nombre: "CoinBase",
		Url:    config.URL,
		Token:  config.Token,
	}
}

func (c Coinbase) Cotizar(simbolos string) (float64, error) {
	url := fmt.Sprintf("%sprices/%s-USD/buy", c.Url, strings.ToUpper(simbolos))
	response, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error al hacer el get a  de %s: %w", url, err)
	}

	var data map[string]map[string]string
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("error al obtener cotizacion de %s: %w", simbolos, err)
	}

	valor, err := strconv.ParseFloat(data["data"]["amount"], 64)
	if err != nil {
		return 0, fmt.Errorf("error parseando el valor de %s: %w", simbolos, err)
	}

	return valor, nil
}

func (c Coinbase) ExisteMoneda(simbolo string) (bool, error) {
	url := c.Url + "currencies/crypto"
	response, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("error al hacer el get a  de %s: %w", url, err)
	}

	// var data map[string][]map[string]string
	var data struct {
		Data []struct {
			Code string `json:"code"`
		}
	}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return false, fmt.Errorf("error al buscar al moneda de %s: %w", simbolo, err)
	}

	for _, moneda := range data.Data {
		if moneda.Code == strings.ToUpper(simbolo) {
			return true, nil
		}
	}

	return false, nil
}

func (c Coinbase) CotizarConcurrente(simbolos []string) ([]ports.Cotizacion, error) {
	ch := make(chan ports.Cotizacion, len(simbolos)+5)
	for _, simbolo := range simbolos {
		go c.cotizarRoutine(simbolo, ch)
	}

	var cotizaciones []ports.Cotizacion
	var errores []error
	for i := 0; i < len(simbolos); i++ {
		cotizacion := <-ch
		if cotizacion.Err != nil {
			errores = append(errores, cotizacion.Err)
			continue
		}
		cotizaciones = append(cotizaciones, cotizacion)
	}

	return cotizaciones, errors.Join(errores...)
}

func (c Coinbase) cotizarRoutine(simbolo string, ch chan ports.Cotizacion) {
	valor, err := c.Cotizar(simbolo)
	if err != nil {
		ch <- ports.Cotizacion{Simbolo: "error", Valor: 0, Err: err}
		return
	}
	ch <- ports.Cotizacion{Simbolo: simbolo, Valor: valor, Err: nil}
}

func (api Coinbase) GetNombre() string {
	return api.Nombre
}
