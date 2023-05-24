/*******************************************************************************
 * YAML config loader and CLI argument parser
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-05-11
 ******************************************************************************/

package main

import (
    "flag"
    "github.com/spf13/viper"
)

func setDefaultConfig() {
    viper.SetDefault("loglevel", 1)
    viper.SetDefault("server.host", "127.0.0.1")
    viper.SetDefault("server.port", 80)
    viper.SetDefault("stats.data", "./")
    viper.SetDefault("stats.log", "./hls.log")
    viper.SetDefault("stats.regex", "/(?P<name>[a-z0-9-_]+)_(?:src|[0-9]+p)/(?P<id>[0-9]+)\\.ts$")
}

func loadConfig() error {
    setDefaultConfig()

    viper.SetConfigName("config")
    viper.SetConfigType("yaml")

    configPath := flag.String("config", "", "Path to config.yaml")
    flag.Parse()
    if *configPath != "" {
        viper.SetConfigFile(*configPath)
    } else {
        viper.AddConfigPath("/etc/hls-utils/")
        viper.AddConfigPath(".")
    }

    return viper.ReadInConfig()
}
