package main

import (
	"fmt"
	"os"

	"grafana-plugin-loader/pkg/config"
	"grafana-plugin-loader/pkg/grafana"
	"grafana-plugin-loader/pkg/utils"
)

func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	appConfig.LogDestination = "./load_grafana.log"

	logger, err := utils.GetLogger(appConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()
	logger.Info("app started")

	pluginData, err := grafana.GetPluginData()
	if err != nil {
		logger.Fatalf("Error: %s", err)
	}
	logger.Infof("downloaded data with %d items", len(pluginData.Items))

	err = grafana.SaveToDB(pluginData, appConfig.DatabaseLocation)
	if err != nil {
		logger.Fatalf("Error: %s", err)
	}
	logger.Info("inserted the data into the database")

	logger.Info("app finished")
}
