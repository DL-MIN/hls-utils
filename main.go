/*******************************************************************************
 * HLS Utils
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2022-11-25
 ******************************************************************************/

package main

import (
	"github.com/spf13/viper"
	"hls-utils/http"
	"hls-utils/logger"
	"hls-utils/stats"
	"hls-utils/terminator"
	"log"
)

func main() {
	if err := loadConfig(); err != nil {
		log.Fatal(err)
	}
	logger.SetLevel(viper.GetInt("loglevel"))

	http.Run()
	stats.Run()
	terminator.WaitGroup.Wait()
}
