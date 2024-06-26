package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	. "hls-utils/logger"
	"hls-utils/terminator"
)

// Run starts a web server on the given IP address and port
func Run() {
	// Set GIN mode
	if Level() == LevelDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()
	router.GET("/auth", authenticate)

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", viper.GetString("server.listen"), viper.GetInt("server.port")),
		Handler:           router,
		ReadHeaderTimeout: time.Minute,
	}

	terminator.WaitGroup.Add(1)
	go func() {
		Infof("start web server and listen to http://%s:%d", viper.GetString("server.listen"), viper.GetInt("server.port"))
		if err := server.ListenAndServe(); err == http.ErrServerClosed {
			Info(err)
		} else {
			Fatal(err)
		}
		terminator.WaitGroup.Done()
	}()

	go func() {
		<-terminator.Signal

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			Warn(err)
		}
	}()
}
