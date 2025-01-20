package main

import (
	"github.com/gin-gonic/gin"
	"hls-utils/controller"
	"hls-utils/types"
)

func setupRouter(router gin.IRouter, config *types.Config, streamManager *types.StreamManager) {
	router.GET("/auth", controller.Authenticate(config))
	router.POST("/auth", controller.Authenticate(config))

	router.GET("/:name/index.m3u8", controller.GetPlaylist(streamManager))
	router.GET("/:name/:client_id/:variant/index.m3u8", controller.GetVariantPlaylist(config, streamManager))
	router.GET("/:name/statistics", controller.GetStatistics(streamManager))
}
