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

	logger, err := utils.GetLogger(appConfig.LogDestination, appConfig.LogMaxSizeBytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger.Println("app started")

	pluginData, err := grafana.GetPluginData()
	if err != nil {
		logger.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	logger.Printf("downloaded data with %d items\n", len(pluginData.Items))

	err = grafana.SaveToDB(pluginData, appConfig.DatabaseLocation)
	if err != nil {
		logger.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	logger.Println("inserted the data into the database")

	logger.Println("app finished")
}
