package main

import (
	"database/sql"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestLogSetup(t *testing.T) {
	tmpDir := t.TempDir()
	logMaxSizeBytes = 1000000
	logDestination = tmpDir + "/app.log"

	err := setupLogging()

	if err != nil {
		t.Errorf("Received error: %s", err)
	}

	logger.Println("test-log-message")

	fileContent, _ := ioutil.ReadFile(logDestination)
	if !strings.Contains(string(fileContent), "test-log-message") {
		t.Errorf("The log destination file does not contain the test log")
	}
}
func TestLogAppend(t *testing.T) {
	tmpDir := t.TempDir()
	logMaxSizeBytes = 1000000
	logDestination = tmpDir + "/app.log"
	os.WriteFile(logDestination, []byte("previous-content"), 0644)

	err := setupLogging()

	if err != nil {
		t.Errorf("Received error: %s", err)
	}

	logger.Println("test-log-message")

	fileContent, _ := ioutil.ReadFile(logDestination)
	if !strings.Contains(string(fileContent), "previous-content") {
		t.Errorf("The log destination file does not contain the previous content")
	}
}

func TestLogRotate(t *testing.T) {
	tmpDir := t.TempDir()
	logMaxSizeBytes = 1
	logDestination = tmpDir + "/app.log"
	os.WriteFile(logDestination, []byte("old-file"), 0644)

	err := setupLogging()

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
	logMaxSizeBytes = 1
	logDestination = tmpDir + "/app.log"
	os.WriteFile(logDestination+".backup", []byte("previous-backup"), 0644)
	os.WriteFile(logDestination, []byte("old-log"), 0644)

	err := setupLogging()

	if err != nil {
		t.Errorf("Received error: %s", err)
	}

	fileContent, _ := ioutil.ReadFile(logDestination + ".backup")
	if strings.Contains(string(fileContent), "previous-backup") {
		t.Errorf("The previous backup content still exists")
	}
}

func TestGetPluginData(t *testing.T) {
	testData := `{
		"other": "stuff",
        "items": [
            {"random": "things", "version": "1", "downloads": 11},
            {"version": "2", "downloads": 22}
        ]
	}`
	expectedPluginData := pluginDataStruct{Items: []pluginDataItem{{"1", 11}, {"2", 22}}}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != pluginDataPath {
			t.Errorf("Expected to request '%s', got: %s", pluginDataPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testData))
	}))
	defer server.Close()

	pluginDataBaseURL = server.URL
	pluginData, err := getPluginData()

	if err != nil {
		t.Errorf("Received error: %s", err)
	}
	if len(pluginData.Items) != 2 {
		t.Errorf("Expected 2 items. Got: %d", len(pluginData.Items))
	}

	if !reflect.DeepEqual(pluginData, expectedPluginData) {
		t.Errorf("Expected %v items. Got: %v", expectedPluginData, pluginData)
	}
}

func TestSaveToDB(t *testing.T) {
	tmpDir := t.TempDir()
	databaseLocation = tmpDir + "/plugin.db"

	db, _ := sql.Open("sqlite", databaseLocation)
	defer db.Close()
	_, _ = db.Exec(
		"CREATE TABLE frser_sqlite (timestamp INTEGER, version TEXT, downloads INTEGER)",
	)

	testPlugindata := pluginDataStruct{Items: []pluginDataItem{{"1", 11}, {"2", 22}}}

	now := time.Now().Unix()
	err := saveToDB(testPlugindata)

	if err != nil {
		t.Errorf("Received error: %s", err)
	}

	rows, _ := db.Query(
		"SELECT timestamp, version, downloads INTEGER FROM frser_sqlite ORDER BY version ASC",
	)
	storedPluginData := pluginDataStruct{Items: []pluginDataItem{}}

	for rows.Next() {
		var timestamp int64
		var item pluginDataItem
		rows.Scan(&timestamp, &item.Version, &item.Downloads)

		storedPluginData.Items = append(storedPluginData.Items, item)
		if math.Abs(float64(now-timestamp)) > 1 {
			t.Errorf("Expected timestamp of roughly %v items. Got: %v", now, timestamp)
		}
	}

	if !reflect.DeepEqual(storedPluginData, testPlugindata) {
		t.Errorf("Expected %v items. Got: %v", testPlugindata, storedPluginData)
	}
}
