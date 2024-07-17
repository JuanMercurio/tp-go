package repos

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/ports"
)

type SentenciaSQL struct {
	Select      string
	From        string
	Where       string
	OrderBy     string
	LimitOffset string
}

func (s SentenciaSQL) toString() string {
	// TODO validaciones
	return s.Select + " " + s.From + " " + s.Where + " " + s.OrderBy + " " + s.LimitOffset
}

func QueryBaseCotizaciones(p ports.ParamCotizaciones) SentenciaSQL {

	var sentencia SentenciaSQL
	sentencia.Select = "SELECT id, id_criptomoneda, fecha, valor, api"
	sentencia.From = "FROM cotizacion"
	sentencia.Where = "WHERE id_criptomoneda IN (" + strings.Join(p.Monedas, ",") + ") AND fecha BETWEEN ? AND ?"
	sentencia.OrderBy = crearSentenciaOrderBy(p.Orden)
	sentencia.LimitOffset = "LIMIT ? OFFSET ?"

	return sentencia
}

func crearSentenciaOrderBy(orden ports.Orden) string {
	var columna, columnaOrden string
	switch orden.TipoOrden {
	case ports.OrdenPorFecha:
		columna = "fecha"
	case ports.OrdenPorNombre:
		columna = "id_criptomoneda"
	case ports.OrdenPorValor:
		columna = "valor"
	}

	if orden.Ascendente {
		columnaOrden = "ASC"
	} else {
		columnaOrden = "DESC"
	}

	return "ORDER BY " + columna + " " + columnaOrden
}

func (r RepositorioMoneda) extraerCotizaciones(rows *sql.Rows) []domain.Cotizacion {
	var cotizaciones []domain.Cotizacion
	var i int
	for rows.Next() {
		i++

		var cotizacion domain.Cotizacion
		var id_moneda int
		var tiempoString string
		if err := rows.Scan(&cotizacion.ID, &id_moneda, &tiempoString, &cotizacion.Valor, &cotizacion.Api); err != nil {
			//TODO no deberia ser aun log fatal
			log.Fatal(err)
		}

		// TODO no ir siempre a la base esto puede ser mas performante si memoizamos o hacemos un join
		moneda, _ := r.BuscarPorId(id_moneda)
		cotizacion.Time, _ = time.Parse("2006-01-02 15:04:05", tiempoString)
		cotizacion.Moneda = moneda

		cotizaciones = append(cotizaciones, cotizacion)
	}

	return cotizaciones
}
