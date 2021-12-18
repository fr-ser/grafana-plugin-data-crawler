package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	// register the SQL driver
	_ "modernc.org/sqlite"
)

var databaseLocationDefault = "./plugin.db"
var databaseLocation string

var pluginDataBaseURL = "https://grafana.com"
var pluginDataPath = "/api/plugins/frser-sqlite-datasource/versions"

var logDestination = "./app.log"
var logMaxSizeBytes int64 = 5000000

type logWriter struct {
	io.Writer
}

// Write appends to the log file and to the stderr
func (w logWriter) Write(b []byte) (n int, err error) {
	fmt.Fprint(os.Stderr, string(append([]byte(time.Now().Format("2006-01-02T15:04:05.999Z")), b...)))
	return w.Writer.Write(append([]byte(time.Now().Format("2006-01-02T15:04:05.999Z")), b...))
}

var logger *log.Logger

type pluginDataItem struct {
	Version   string `json:"version"`
	Downloads int32  `json:"Downloads"`
}

type pluginDataStruct struct {
	Items []pluginDataItem `json:"items"`
}

func setupLogging() (err error) {
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

	file, err := os.OpenFile(logDestination, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("error opening the log file: %s", err)
	}

	logger = log.New(&logWriter{file}, " | ", 0)
	return nil
}

func getPluginData() (pluginData pluginDataStruct, err error) {
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

func saveToDB(pluginData pluginDataStruct) (err error) {
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

func main() {
	if value, ok := os.LookupEnv("DB_LOCATION"); ok {
		databaseLocation = value
	}
	databaseLocation = databaseLocationDefault

	err := setupLogging()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger.Println("app started")

	pluginData, err := getPluginData()
	if err != nil {
		logger.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	logger.Printf("downloaded data with %d items\n", len(pluginData.Items))

	err = saveToDB(pluginData)
	if err != nil {
		logger.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	logger.Println("inserted the data into the database")

	logger.Println("app finished")
}
