package utils

import (
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"grafana-plugin-loader/pkg/config"
)

// maybeRotateFile checks the log file and rotates it if necessary
// it does so by creating a file appended with .backup
func maybeRotateFile(logDestination string, logMaxSizeBytes int64) error {
	var logSize int64

	fi, err := os.Stat(logDestination)
	if err == nil {
		logSize = fi.Size()
	}

	if logSize > logMaxSizeBytes {
		err = os.Remove(logDestination + ".backup")
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("error removing the previous log backup: %s", err)
		}
		err = os.Rename(logDestination, logDestination+".backup")
		if err != nil {
			return fmt.Errorf("error creating the new log backup: %s", err)
		}
	}
	return nil
}

// GetLogger returns a fully set up logger
func GetLogger(appConfig config.AppConfig) (logger *zap.SugaredLogger, err error) {
	err = maybeRotateFile(appConfig.LogDestination, appConfig.LogMaxSizeBytes)
	if err != nil {
		return nil, err
	}

	var logConfig = zap.NewProductionConfig()

	logConfig.Level.UnmarshalText([]byte(appConfig.LogLevel))
	logConfig.OutputPaths = []string{"stderr", appConfig.LogDestination}
	logConfig.Encoding = "console"
	logConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logConfig.DisableStacktrace = true

	rawLogger, err := logConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("error creating logger from config: %v", err)
	}

	return rawLogger.Sugar(), nil
}
