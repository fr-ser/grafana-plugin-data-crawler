package utils

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLogRotate(t *testing.T) {
	tmpDir := t.TempDir()
	var logMaxSizeBytes int64 = 1
	logDestination := tmpDir + "/app.log"
	_ = os.WriteFile(logDestination, []byte("old-file"), 0644)

	err := maybeRotateFile(logDestination, logMaxSizeBytes)

	if err != nil {
		t.Errorf("Received error: %s", err)
	}

	fileContent, _ := ioutil.ReadFile(logDestination + ".backup")
	if !strings.Contains(string(fileContent), "old-file") {
		t.Errorf("The newly created backup file does not have the right content")
	}
}

func TestLogRotateDeletePrevious(t *testing.T) {
	tmpDir := t.TempDir()
	var logMaxSizeBytes int64 = 1
	logDestination := tmpDir + "/app.log"
	_ = os.WriteFile(logDestination+".backup", []byte("previous-backup"), 0644)
	_ = os.WriteFile(logDestination, []byte("old-log"), 0644)

	err := maybeRotateFile(logDestination, logMaxSizeBytes)

	if err != nil {
		t.Errorf("Received error: %s", err)
	}

	fileContent, _ := ioutil.ReadFile(logDestination + ".backup")
	if strings.Contains(string(fileContent), "previous-backup") {
		t.Errorf("The previous backup content still exists")
	}
}
