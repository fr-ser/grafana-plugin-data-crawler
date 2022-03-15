package grafana

import (
	"database/sql"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

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
		_, _ = w.Write([]byte(testData))
	}))
	defer server.Close()

	pluginDataBaseURL = server.URL
	pluginData, err := GetPluginData()

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
	databaseLocation := tmpDir + "/plugin.db"

	db, _ := sql.Open("sqlite", databaseLocation)
	defer db.Close()
	_, _ = db.Exec(
		"CREATE TABLE frser_sqlite (timestamp INTEGER, version TEXT, downloads INTEGER)",
	)

	testPlugindata := pluginDataStruct{Items: []pluginDataItem{{"1", 11}, {"2", 22}}}

	now := time.Now().Unix()
	err := SaveToDB(testPlugindata, databaseLocation)

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
		_ = rows.Scan(&timestamp, &item.Version, &item.Downloads)

		storedPluginData.Items = append(storedPluginData.Items, item)
		if math.Abs(float64(now-timestamp)) > 1 {
			t.Errorf("Expected timestamp of roughly %v items. Got: %v", now, timestamp)
		}
	}

	if !reflect.DeepEqual(storedPluginData, testPlugindata) {
		t.Errorf("Expected %v items. Got: %v", testPlugindata, storedPluginData)
	}
}
