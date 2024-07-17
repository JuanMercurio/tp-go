package paprika

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/juanmercurio/tp-go/internal/adapters/config"
)

type Paprika struct {
	Nombre string
	Url    string
	Token  string
	MapIds map[string]string
}

func Crear(config *config.APIConfig) Paprika {
	return Paprika{
		Nombre: "Paprika",
		Url:    config.URL,
		Token:  config.Token,
		// problemas con la concurrencia falta ponerle un mutex o cambiar el disenio
		MapIds: make(map[string]string),
	}
}

func (api Paprika) GetNombre() string {
	return api.Nombre
}

func (api Paprika) ExisteMoneda(simbolo string) (bool, error) {
	id, err := api.obtenerId(strings.ToUpper(simbolo))
	if err != nil {
		return false, err
	}

	if id == "" {
		return false, nil
	}

	return true, nil

}

func (api Paprika) Cotizar(simboloMoneda string) (float64, error) {
	id, ok := api.MapIds[simboloMoneda]
	if !ok {
		var err error
		id, err = api.obtenerId(simboloMoneda)
		if err != nil {
			return 0, err
		}
	}

	url := fmt.Sprintf("%stickers/%s", api.Url, id)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("no pude obtener la cotizacion de paprika: %v", err)
	}

	var cotizacion struct {
		Quotes struct {
			USD struct {
				Price float64 `json:"price"`
			}
		} `json:"quotes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&cotizacion); err != nil {
		return 0, fmt.Errorf("error en la deodificacion del precio")
	}

	return cotizacion.Quotes.USD.Price, nil
}

func (p Paprika) obtenerId(simbolo string) (string, error) {

	var monedas []struct {
		Id     string `json:"id"`
		Symbol string `json:"symbol"`
	}

	resp, err := http.Get(p.Url + "coins")
	if err != nil {
		return "", fmt.Errorf("no pude obtener las monedas de paprika: %v", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&monedas); err != nil {
		return "", fmt.Errorf("no pude decodificar las monedas de paprika: %v", err)
	}

	for _, moneda := range monedas {
		if moneda.Symbol == simbolo {
			return moneda.Id, nil
		}
	}
	return "", nil
}

// func (api Paprika) ObtenerCotizacion(moneda string, codigo string) (cotizador.Cotizacion, error) {

// 	url := fmt.Sprintf("%s/%s", api.Url, "coins")
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return cotizador.Cotizacion{}, err
// 	}

// 	req.Header.Add("Authorization", "51e0e631-1f99-46d1-84c0-44667f8070fa")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return cotizador.Cotizacion{}, err
// 	}
// 	defer resp.Body.Close()

// 	var m monedas
// 	err = json.NewDecoder(resp.Body).Decode(&m)
// 	if err != nil {
// 		return cotizador.Cotizacion{}, err
// 	}

// 	fmt.Println(m)

// 	return cotizador.Cotizacion{}, nil
// }
