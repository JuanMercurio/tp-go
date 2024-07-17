package repos

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/juanmercurio/tp-go/internal/core/domain"
	"github.com/juanmercurio/tp-go/internal/ports"
)

type GetCotizacionesSQL struct {
	Monedas     string
	FechaInicio time.Time
	FechaFinal  time.Time
	Cant        int
	Offset      int
}

func aSQL(p ports.ParamCotizaciones) GetCotizacionesSQL {
	var parametroSQL GetCotizacionesSQL

	parametroSQL.Monedas = strings.Join(p.Monedas, ",")
	parametroSQL.FechaInicio = p.FechaInicial
	parametroSQL.FechaFinal = p.FechaFinal
	parametroSQL.Cant = p.CantPaginas * p.TamPaginas
	parametroSQL.Offset = p.PaginaInicial

	return parametroSQL
}

func queryBaseCotizaciones(p ports.ParamCotizaciones) string {
	selectClause := crearSentenciaSelect("cotizacion", "id", "id_criptomoneda", "fecha", "valor")
	orderBy := crearSentenciaOrderBy(p.Orden)
	where := "WHERE id_criptomoneda IN (?) AND fecha BETWEEN ? AND ?"
	limitOffset := "LIMIT ? OFFSET ?"
	return selectClause + " " + where + " " + orderBy + " " + limitOffset
}

func crearSentenciaSelect(tabla string, columnas ...string) string {
	return "SELECT " + strings.Join(columnas, ",") + " FROM " + tabla
}

func crearSentenciaOrderBy(orden ports.Orden) string {
	var columna, columnaOrden string
	switch orden.TipoOrden {
	case ports.OrdenPorFecha:
		columna = "fecha"
	case ports.OrdenPorNombre:
		columna = "id"
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
	for rows.Next() {
		var cotizacion domain.Cotizacion
		var id_moneda int
		var tiempoString string
		if err := rows.Scan(&cotizacion.ID, &id_moneda, &tiempoString, &cotizacion.Valor); err != nil {
			//TODO no deberia ser aun log fatal
			log.Fatal(err)
		}

		// TODO no ir siempre a la base esto puede ser mas performante si memoizamos
		moneda, _ := r.BuscarPorId(id_moneda)
		cotizacion.Time, _ = time.Parse("2006-01-02 15:04:05", tiempoString)
		cotizacion.Moneda = moneda

		cotizaciones = append(cotizaciones, cotizacion)
	}

	return cotizaciones
}
