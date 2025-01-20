package main

import (
	"github.com/spacecafe/gobox/httpserver"
	"github.com/spacecafe/gobox/logger"
	"github.com/spacecafe/gobox/terminator"
	"hls-utils/types"
)

func main() {
	term := terminator.New(terminator.NewConfig())
	config := types.NewConfig()

	streamManager := types.NewStreamManager(config)
	if err := streamManager.Start(term.FullTracking()); err != nil {
		logger.Fatal(err)
	}

	server := httpserver.NewHTTPServer(config.HTTPServer)
	setupRouter(server.Router, config, streamManager)
	if err := server.Start(term.FullTracking()); err != nil {
		logger.Fatal(err)
	}

	term.Wait()
}
