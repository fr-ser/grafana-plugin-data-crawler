package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"

	"grafana-plugin-loader/pkg/grafana"
	"grafana-plugin-loader/pkg/utils"
)

type AppConfig struct {
	DatabaseLocation string `env:"DB_LOCATION" envDefault:"./plugin.db"`

	LogDestination  string `env:"LOG_DESTINATION" envDefault:"./app.log"`
	LogMaxSizeBytes int64  `env:"LOG_MAX_SIZE_BYTES" envDefault:"5000000"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info"`
}

func loadConfig() (config AppConfig, err error) {
	if err := env.Parse(&config); err != nil {
		return config, fmt.Errorf("error loading configuration: %v", err)
	}

	return config, nil
}

func main() {
	appConfig, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	appConfig.LogDestination = "./load_grafana.log"

	logger, err := utils.GetLogger(appConfig.LogDestination, appConfig.LogMaxSizeBytes, appConfig.LogLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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

	_ = logger.Sync()
}
