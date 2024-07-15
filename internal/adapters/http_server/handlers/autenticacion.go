package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func autenticarUsuario(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	if token == "" {
		return errors.New("se necesita un token para realizar esta operacion")
	}

	//TODO agregar jwt aca
	if token != "token" {
		return errors.New("token incorrecto")
	}

	return nil
}
