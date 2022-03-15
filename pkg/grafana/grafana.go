package grafana

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	// register the SQL driver
	_ "modernc.org/sqlite"
)

var pluginDataBaseURL = "https://grafana.com"
var pluginDataPath = "/api/plugins/frser-sqlite-datasource/versions"

type pluginDataItem struct {
	Version   string `json:"version"`
	Downloads int32  `json:"Downloads"`
}

type pluginDataStruct struct {
	Items []pluginDataItem `json:"items"`
}

// GetPluginData retrieves the data about the sqlite plugin from the Grafana API
func GetPluginData() (pluginData pluginDataStruct, err error) {
	httpClient := http.Client{Timeout: time.Second * 30}
	req, err := http.NewRequest(http.MethodGet, pluginDataBaseURL+pluginDataPath, nil)
	if err != nil {
		return pluginData, fmt.Errorf("error creating the request: %s", err)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return pluginData, fmt.Errorf("error doing the request: %s", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return pluginData, fmt.Errorf("error reading the request body: %s", err)
	}

	err = json.Unmarshal(body, &pluginData)
	if err != nil {
		return pluginData, fmt.Errorf("error parsing the request body: %s", err)
	}

	return pluginData, nil
}

// SaveToDB stores the pluginData in an SQLite database
func SaveToDB(pluginData pluginDataStruct, databaseLocation string) (err error) {
	db, err := sql.Open("sqlite", databaseLocation)
	if err != nil {
		return fmt.Errorf("error opending the database: %s", err)
	}
	defer db.Close()

	for _, item := range pluginData.Items {
		_, err = db.Exec(
			"INSERT INTO frser_sqlite (timestamp, version, downloads) VALUES (?, ?, ?)",
			time.Now().Unix(), item.Version, item.Downloads,
		)
		if err != nil {
			return fmt.Errorf("error inserting row: %s", err)
		}
	}

	return nil
}
