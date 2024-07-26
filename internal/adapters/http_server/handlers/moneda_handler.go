package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanmercurio/tp-go/internal/ports"
)

type MonedaHandler struct {
	srv ports.ServicioMonedas
}

func CrearHandlerMoneda(srv ports.ServicioMonedas) *MonedaHandler {
	return &MonedaHandler{
		srv: srv,
	}
}

// @Summary		Busca todas las monedas
// @Description	Obtiene una lista de todos las monedas disponibles.
// @Tags			Moneda
// @Accept			json
// @Produce		json
// @Success		200	{object}	[]ports.MonedaOutputDTO
// @Router			/monedas [get]
func (mh MonedaHandler) BuscarMonedas(c *gin.Context) {
	todos, err := mh.srv.BuscarTodos()
	if err != nil {
		c.JSON(http.StatusConflict, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

// @Summary		Da de alta una moneda
// @Description	Si tenemos las credenciales podemos dar de alta una moneda
// @Tags			Moneda
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string	true	"Token de autorización"
// @Param			simbolo			query		string	true	"Simbolo de la moneda"
// @Param			nombre			query		string	false	"Nombre de la moneda nueva"
// @Success		200				{object}	int
// @Failure		400				{object}	error
// @Router			/monedas [post]
func (mh MonedaHandler) AltaMoneda(c *gin.Context) {
	nombre := c.Query("nombre")
	simbolo := c.Query("simbolo")

	id, err := mh.srv.AltaMoneda(nombre, simbolo)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// mh.srv.CotizarNuevaMoneda(id, simbolo)

	c.JSON(http.StatusOK, gin.H{"id": id})
}
