package cotizador

// import (
// 	"fmt"
// )

// type GestorApisCotizadoras struct {
// 	cotizadores map[string]APICotizadora
// }

// func Crear(apis ...APICotizadora) *GestorApisCotizadoras {

// 	var cotizadores = make(map[string]APICotizadora)
// 	for _, api := range apis {
// 		cotizadores[api.GetNombre()] = api
// 	}

// 	return &GestorApisCotizadoras{
// 		cotizadores: cotizadores,
// 	}
// }

// func (c GestorApisCotizadoras) Cotizar(api, simbolo string) (float64, error) {
// 	cotizador, ok := c.cotizadores[api]
// 	if !ok {
// 		return 0, fmt.Errorf("no se encontro el cotizador %s, para la moneda: %s", api, simbolo)
// 	}

// 	return cotizador.Cotizar(simbolo)
// }

// // si no existe en algun api la moneda se da de baja
// func (c GestorApisCotizadoras) ExisteMoneda(simbolo string) (bool, error) {
// 	for _, cotizador := range c.cotizadores {
// 		existe, err := cotizador.ExisteMoneda(simbolo)
// 		if err != nil {
// 			return false, fmt.Errorf("la moneda no existe en el servicio %s: %w", cotizador.GetNombre(), err)
// 		}

// 		if !existe {
// 			return false, nil
// 		}
// 	}

// 	return true, nil
// }

// func (c GestorApisCotizadoras) GetNombre() string {
// 	return "GestorApisCotizadoras"
// }
