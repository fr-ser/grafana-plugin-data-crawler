package utils

import (
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
func GetLogger(logDestination string, logMaxSizeBytes int64, logLevel string) (logger *zap.SugaredLogger, err error) {
	err = maybeRotateFile(logDestination, logMaxSizeBytes)
	if err != nil {
		return nil, err
	}

	var logConfig = zap.NewProductionConfig()

	err = logConfig.Level.UnmarshalText([]byte(logLevel))
	if err != nil {
		return nil, fmt.Errorf("error setting log level: %v", err)
	}
	logConfig.OutputPaths = []string{"stderr", logDestination}
	logConfig.Encoding = "console"
	logConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logConfig.DisableStacktrace = true

	rawLogger, err := logConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("error creating logger from config: %v", err)
	}

	return rawLogger.Sugar(), nil
}
