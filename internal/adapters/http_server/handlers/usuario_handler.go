package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
// @Tags		Usuarios
// @Accept		json
// @Produce	json
// @Param		username			query		string	true	"nombre de usuario elegido"
// @Param		nombre				query		string	true	"nombre del nuevo usuario"
// @Param		apellido			query		string	true	"apellido del nuevo usuario"
// @Param		email				query		string	true	"email del nuevo usuario"
// @Param		tipo_documento		query		string	true	"dni | cedula | pasaporte"
// @Param		documento_numero	query		string	true	"identificador de documento elegido"
// @Param		fecha_nacimiento	query		string	true	"formato YYYY-MM-DD HH:MM:SS"
// @Success	200					{object}	string
// @Failure	400					{object}	string
// @Router		/usuarios [post]
func (mh UsuarioHandler) AltaUsuario(c *gin.Context) {
	parametros, err := ValidarParamsAltaUsuario(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := mh.srv.AltaUsuario(parametros)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary	Da de baja a un usuario
// @Tags		Usuarios
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
// @Tags		Usuarios
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

func ValidarParamsAltaUsuario(c *gin.Context) (ports.AltaUsuarioParams, error) {

	if c.Query("nombre") == "" {
		return ports.AltaUsuarioParams{}, errors.New("debe proporcionar un nombre de usuario")
	}

	if c.Query("email") == "" {
		return ports.AltaUsuarioParams{}, errors.New("debe proporcionar un email usuario")
	}

	if c.Query("apellido") == "" {
		return ports.AltaUsuarioParams{}, errors.New("debe proporcionar un apellido")
	}

	if c.Query("fecha_nacimiento") == "" {
		return ports.AltaUsuarioParams{}, errors.New("debe proporcionar una fecha de nacimiento")
	}

	if c.Query("tipo_documento") == "" {
		return ports.AltaUsuarioParams{}, errors.New("debe proporcionar un tipo de documento")
	}

	if c.Query("documento_numero") == "" {
		return ports.AltaUsuarioParams{}, errors.New("debe proporcionar un numero de documento")
	}

	if c.Query("username") == "" {
		return ports.AltaUsuarioParams{}, errors.New("debe proporcionar un username")
	}

	fecha, err := stringAFecha(c.Query("fecha_nacimiento"))
	if err != nil {
		return ports.AltaUsuarioParams{}, err
	}

	if esMenor(fecha) {
		return ports.AltaUsuarioParams{}, errors.New("debes ser mayor de edad para usar esta plataforma")
	}

	tipoDoc := strings.ToUpper(c.Query("tipo_documento"))
	if tipoDocValido(tipoDoc) {
		return ports.AltaUsuarioParams{}, errors.New("tipo invalido de documento")
	}

	parametros := ports.AltaUsuarioParams{
		Username:          c.Query("username"),
		Nombre:            c.Query("nombre"),
		Apellido:          c.Query("apellido"),
		FechaDeNacimiento: fecha,
		Documento:         c.Query("documento_numero"),
		TipoDocumento:     tipoDoc,
		Email:             c.Query("email"),
	}

	return parametros, nil
}

// @Summary		Actualizar atributos de usuario
// @Description	Actualiza parcialmente un usuario por su ID
// @Tags			Usuarios
// @Accept			json
// @Param			id		path		integer			true	"ID del usuario a actualizar"
// @Param			body	body		[]ports.Patch	true	"Datos de actualización en JSON"
// @Success		200		{object}	domain.Usuario
// @Failure		400		{object}	string
// @Router			/usuarios/{id} [patch]
func (uh UsuarioHandler) ActualizarUsuario(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var cambios []ports.Patch
	if err := c.BindJSON(&cambios); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("formato incorrecto: %w", err).Error()})
		return
	}

	if err := validarActualizarUsuario(cambios); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usuario, err := uh.srv.PatchUsuario(id, cambios)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, usuario)
}

func esMenor(fecha time.Time) bool {
	return time.Since(fecha).Hours()/24/365 < 18
}

func validarActualizarUsuario(cambios []ports.Patch) error {
	for _, c := range cambios {

		if c.Path == "fecha" {
			fecha, err := time.Parse("2006-01-02 15:04:05", c.NuevoValor.(string))
			if err != nil {
				return err
			}
			if esMenor(fecha) {
				return errors.New("debe ser mayor de edad")
			}
		}

		if c.Path == "tipo_doc" {
			if !tipoDocValido(c.NuevoValor.(string)) {
				return errors.New("el tipo de doc debe ser DNI | PASAPORTE | CEDULA")
			}

		}
	}

	return nil
}

func tipoDocValido(tipo string) bool {
	return tipo != "DNI" && tipo != "PASAPORTE" && tipo != "CEDULA"
}
