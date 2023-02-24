package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping() func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}
