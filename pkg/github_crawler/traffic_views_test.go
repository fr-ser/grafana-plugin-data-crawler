package github_crawler

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/v43/github"
)

func intPointer(x int) *int {
	return &x
}

func TestGetAndStoreTrafficViews(t *testing.T) {
	tmpDir := t.TempDir()
	databaseLocation := tmpDir + "/plugin.db"

	db, _ := sql.Open("sqlite", databaseLocation)
	defer db.Close()
	_, _ = db.Exec(
		"CREATE TABLE github_traffic_views (timestamp INTEGER PRIMARY KEY, count INTEGER, uniques INTEGER)",
	)

	testViews := []*github.TrafficData{
		{Count: intPointer(33), Uniques: intPointer(22), Timestamp: &github.Timestamp{Time: time.Unix(123456789, 0)}},
		{Count: intPointer(55), Uniques: intPointer(44), Timestamp: &github.Timestamp{Time: time.Unix(223456789, 0)}},
	}

	err := StoreTrafficViews(testViews, databaseLocation)

	if err != nil {
		t.Errorf("Received error: %s", err)
	}

	rows, _ := db.Query(
		"SELECT timestamp, count, uniques FROM github_traffic_views ORDER BY timestamp ASC",
	)
	storedViews := []*github.TrafficData{}

	for rows.Next() {
		var timestamp int64
		item := github.TrafficData{}
		_ = rows.Scan(&timestamp, &item.Count, &item.Uniques)

		item.Timestamp = &github.Timestamp{Time: time.Unix(timestamp, 0)}

		storedViews = append(storedViews, &item)
	}

	if !reflect.DeepEqual(storedViews, testViews) {
		t.Errorf("Expected %+v items. Got: %+v", testViews, storedViews)
	}
}

func TestGetAndStoreTrafficViewsUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	databaseLocation := tmpDir + "/plugin.db"

	db, _ := sql.Open("sqlite", databaseLocation)
	defer db.Close()
	_, _ = db.Exec(
		"CREATE TABLE github_traffic_views (timestamp INTEGER PRIMARY KEY, count INTEGER, uniques INTEGER)",
	)
	_, _ = db.Exec(
		"INSERT INTO github_traffic_views (timestamp, count, uniques) VALUES (123, 2, 1)",
	)

	testViews := []*github.TrafficData{
		{Count: intPointer(33), Uniques: intPointer(22), Timestamp: &github.Timestamp{Time: time.Unix(123, 0)}},
		{Count: intPointer(55), Uniques: intPointer(44), Timestamp: &github.Timestamp{Time: time.Unix(456, 0)}},
	}

	err := StoreTrafficViews(testViews, databaseLocation)

	if err != nil {
		t.Errorf("Received error: %s", err)
	}

	rows, _ := db.Query(
		"SELECT timestamp, count, uniques FROM github_traffic_views ORDER BY timestamp ASC",
	)
	storedViews := []*github.TrafficData{}

	for rows.Next() {
		var timestamp int64
		item := github.TrafficData{}
		_ = rows.Scan(&timestamp, &item.Count, &item.Uniques)

		item.Timestamp = &github.Timestamp{Time: time.Unix(timestamp, 0)}

		storedViews = append(storedViews, &item)
	}

	if !reflect.DeepEqual(storedViews, testViews) {
		t.Errorf("Expected %+v items. Got: %+v", testViews, storedViews)
	}
}
