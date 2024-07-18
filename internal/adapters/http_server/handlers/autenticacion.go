package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Autenticar(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "se necesita token para realizar esta operacion"})
		c.Abort()
		return
	}

	if token != "token" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token incorrecto"})
		c.Abort()
		return
	}

	c.Next()
}
