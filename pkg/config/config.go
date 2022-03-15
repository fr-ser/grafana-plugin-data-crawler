package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type AppConfig struct {
	DatabaseLocation string `env:"DB_LOCATION" envDefault:"./plugin.db"`

	LogDestination  string `env:"LOG_DESTINATION" envDefault:"./app.log"`
	LogMaxSizeBytes int64  `env:"LOG_MAX_SIZE_BYTES" envDefault:"5000000"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info"`
}

// LoadConfig reads the environment and returns a configuration
func LoadConfig() (config AppConfig, err error) {
	if err := env.Parse(&config); err != nil {
		return config, fmt.Errorf("error loading configuration: %v", err)
	}

	return config, nil
}
