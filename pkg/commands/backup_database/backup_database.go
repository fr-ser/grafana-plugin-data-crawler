package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v6"

	"grafana-plugin-loader/pkg/backup"
	"grafana-plugin-loader/pkg/utils"
)

type AppConfig struct {
	DatabaseLocation string `env:"DB_LOCATION" envDefault:"./plugin.db"`

	LogDestination  string `env:"LOG_DESTINATION" envDefault:"./backup.log"`
	LogMaxSizeBytes int64  `env:"LOG_MAX_SIZE_BYTES" envDefault:"5000000"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info"`

	DropboxToken string `env:"DROPBOX_TOKEN,file" envDefault:"./.dropbox_token"`
}

func loadConfig() (config AppConfig, err error) {
	if err := env.Parse(&config); err != nil {
		return config, fmt.Errorf("error loading configuration: %v", err)
	}
	config.DropboxToken = strings.TrimSuffix(config.DropboxToken, "\n")

	return config, nil
}
func main() {

	appConfig, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger, err := utils.GetLogger(appConfig.LogDestination, appConfig.LogMaxSizeBytes, appConfig.LogLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger.Info("app started")

	err = backup.CheckDatabaseConsistency(appConfig.DatabaseLocation)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("database is valid")

	err = backup.UploadDatabase(appConfig.DatabaseLocation, appConfig.DropboxToken)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("database uploaded")

	_ = logger.Sync()
}
