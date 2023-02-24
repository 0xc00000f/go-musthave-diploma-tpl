package handlers

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", Ping())

	return r
}
