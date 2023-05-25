package http

import "github.com/gin-gonic/gin"

var router *gin.Engine

func init() {
	router = gin.Default()
	router.GET("/auth", authenticate)
}
