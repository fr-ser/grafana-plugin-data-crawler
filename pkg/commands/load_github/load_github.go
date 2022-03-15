package main

import (
	"fmt"
	"grafana-plugin-loader/pkg/config"
	"grafana-plugin-loader/pkg/utils"
	"os"
)

func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	appConfig.LogDestination = "./load_github.log"

	logger, err := utils.GetLogger(appConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()
	logger.Info("app started")
}
