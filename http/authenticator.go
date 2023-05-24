package http

import (
    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
    "net/http"
)

func authenticate(c *gin.Context) {
    if c.Request.FormValue("call") != "publish" {
        c.AbortWithStatus(http.StatusMethodNotAllowed)
    }

    configPath := "streams." + c.Request.FormValue("name")

    if viper.IsSet(configPath) && viper.GetString(configPath) == c.Request.FormValue("auth") {
        c.AbortWithStatus(http.StatusOK)
    }

    c.AbortWithStatus(http.StatusUnauthorized)
}
