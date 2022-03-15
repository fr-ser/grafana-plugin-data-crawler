package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type logWriter struct {
	io.Writer
}

// Write appends to the log file and to the stderr
func (w logWriter) Write(b []byte) (n int, err error) {
	fmt.Fprint(os.Stderr, string(append([]byte(time.Now().Format("2006-01-02T15:04:05.999Z")), b...)))
	return w.Writer.Write(append([]byte(time.Now().Format("2006-01-02T15:04:05.999Z")), b...))
}

// GetLogger returns a fully set up logger
func GetLogger(logDestination string, logMaxSizeBytes int64) (logger *log.Logger, err error) {
	var logSize int64

	fi, err := os.Stat(logDestination)
	if err == nil {
		logSize = fi.Size()
	}

	if logSize > logMaxSizeBytes {
		err = os.Remove(logDestination + ".backup")
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("error removing the previous log backup: %s", err)
		}
		err = os.Rename(logDestination, logDestination+".backup")
		if err != nil {
			return nil, fmt.Errorf("error creating the new log backup: %s", err)
		}
	}

	file, err := os.OpenFile(logDestination, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening the log file: %s", err)
	}

	return log.New(&logWriter{file}, " | ", 0), nil
}
