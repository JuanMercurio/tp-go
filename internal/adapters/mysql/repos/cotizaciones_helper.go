package repos

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/ports"
)

func QueryBaseCotizaciones(p ports.Filter) queryBuilder {

	var sentencia queryBuilder
	sentencia.From = "FROM cotizacion"

	if len(p.Monedas) > 0 {
		monedas := strings.Join(p.Monedas, ",")
		sentencia.AddWhere("id_criptomoneda IN (" + monedas + ")")
	}

	if !p.FechaInicial.IsZero() {
		fechaSQL := p.FechaInicial.Format("2006-01-02 15:04:05")
		sentencia.AddWhere("fecha >= " + "'" + fechaSQL + "'")
	}

	if !p.FechaFinal.IsZero() {
		fechaSQL := p.FechaFinal.Format("2006-01-02 15:04:05")
		sentencia.AddWhere("fecha <= " + "'" + fechaSQL + "'")
	}

	if p.Usuario != 0 {
		sentencia.AddJoin("usuario_criptomoneda", "cotizacion.id_criptomoneda = usuario_criptomoneda.id_criptomoneda")
		sentencia.AddWhere(fmt.Sprintf("usuario_criptomoneda.id_usuario = %d", (p.Usuario)))
	}

	sentencia.AddOrderBy(ordenToString(p.Orden))

	return sentencia
}

func ordenToString(orden ports.Orden) string {
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

	return columna + " " + columnaOrden
}

func (r RepositorioMoneda) extraerCotizaciones(rows *sql.Rows) []domain.Cotizacion {
	var cotizaciones []domain.Cotizacion
	var i int
	for rows.Next() {
		i++

		var cotizacion domain.Cotizacion
		var id_moneda int
		var tiempoString string
		if err := rows.Scan(&id_moneda, &tiempoString, &cotizacion.Valor, &cotizacion.Api); err != nil {
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
