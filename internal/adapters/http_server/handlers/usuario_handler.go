package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juanmercurio/tp-go/internal/ports"
)

type UsuarioHandler struct {
	srv ports.ServicioUsuarios
}

func CrearHandlerUsuario(srv ports.ServicioUsuarios) *UsuarioHandler {
	return &UsuarioHandler{
		srv: srv,
	}
}

// @Summary	Crear un usuario
// @Tags		Moneda
// @Accept		json
// @Produce	json
// @Param		nombre	query		string	true	"nombre del nuevo usuario"
// @Success	200		{object}	string
// @Failure	400		{object}	string
// @Router		/usuarios [post]
func (mh UsuarioHandler) AltaUsuario(c *gin.Context) {
	nombre := c.Query("nombre")
	if nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe proporcionar un nombre de usuario"})
		return
	}

	id, err := mh.srv.AltaUsuario(nombre)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary	Da de baja a un usuario
// @Tags		Moneda
// @Accept		json
// @Produce	json
// @Param		id	query		int	true	"id del usuario que se desea dar de baja"
// @Success	200	{object}	string
// @Failure	400	{object}	string
// @Router		/usuarios [delete]
func (mh UsuarioHandler) BajaUsuario(c *gin.Context) {

	idString := c.Query("id")
	if idString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe proporcionar un nombre de usuario"})
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el id debe ser un numero"})
	}

	err = mh.srv.BajaUsuario(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "El usuario se elimino correctamente")
}

// @Summary	Lista de usuario registrados
// @Tags		Moneda
// @Accept		json
// @Produce	json
// @Success	200	{object}	[]domain.Usuario
// @Failure	400	{object}	string
// @Router		/usuarios [get]
func (mh UsuarioHandler) BuscarUsuarios(c *gin.Context) {

	usuarios, err := mh.srv.BuscarTodos()
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, usuarios)
}
