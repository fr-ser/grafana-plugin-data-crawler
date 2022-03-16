package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"

	"grafana-plugin-loader/pkg/github_crawler"
	"grafana-plugin-loader/pkg/utils"
)

type AppConfig struct {
	DatabaseLocation string `env:"DB_LOCATION" envDefault:"./plugin.db"`

	LogDestination  string `env:"LOG_DESTINATION" envDefault:"./app.log"`
	LogMaxSizeBytes int64  `env:"LOG_MAX_SIZE_BYTES" envDefault:"5000000"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info"`

	GithubToken string `env:"GITHUB_TOKEN,file" envDefault:"./.github_token"`
}

func loadConfig() (config AppConfig, err error) {
	if err := env.Parse(&config); err != nil {
		return config, fmt.Errorf("error loading configuration: %v", err)
	}

	config.GithubToken = strings.TrimSuffix(config.GithubToken, "\n")

	return config, nil
}

func main() {
	appConfig, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	appConfig.LogDestination = "./load_github.log"

	logger, err := utils.GetLogger(appConfig.LogDestination, appConfig.LogMaxSizeBytes, appConfig.LogLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger.Info("app started")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: appConfig.GithubToken})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	logger.Info("Github Client created")

	pages, _, err := client.Repositories.ListTrafficViews(ctx, "fr-ser", "grafana-sqlite-datasource", nil)
	if err != nil {
		logger.Fatalf("error during traffic view retrieval: %v", err)
	}
	logger.Info("Traffic Views downloaded")

	err = github_crawler.StoreTrafficViews(pages.Views, appConfig.DatabaseLocation)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	logger.Info("Traffic Views stored")

	_ = logger.Sync()
}
