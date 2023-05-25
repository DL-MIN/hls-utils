package http

import "github.com/gin-gonic/gin"

// router is the global routing engine
var router *gin.Engine

// init will create the routes
func init() {
	router = gin.Default()
	router.GET("/auth", authenticate)
}
