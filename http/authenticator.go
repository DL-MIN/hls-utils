package http

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

// authenticate returns HTTP_OK if the given authentication token matches the stored one.
// Otherwise, it returns HTTP_UNAUTHORIZED.
func authenticate(c *gin.Context) {
	if c.Request.FormValue("call") != "publish" {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}

	configPath := "streams." + c.Request.FormValue("name")

	if viper.IsSet(configPath) && viper.GetString(configPath) == c.Request.FormValue("auth") {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}
