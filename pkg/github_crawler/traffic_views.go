package github_crawler

import (
	"database/sql"
	"fmt"

	"github.com/google/go-github/v43/github"

	// register the SQL driver
	_ "modernc.org/sqlite"
)

// StoreTrafficViews stores page views from Github in the SQLite database
func StoreTrafficViews(views []*github.TrafficData, databaseLocation string) error {
	db, err := sql.Open("sqlite", databaseLocation)
	if err != nil {
		return fmt.Errorf("error opending the database: %s", err)
	}
	defer db.Close()

	for _, item := range views {
		_, err = db.Exec(
			`INSERT INTO github_traffic_views (timestamp, count, uniques) VALUES (?, ?, ?)
			ON CONFLICT(timestamp) DO UPDATE SET count=excluded.count, uniques = excluded.uniques`,
			item.Timestamp.Unix(), item.Count, item.Uniques,
		)
		if err != nil {
			return fmt.Errorf("error inserting row: %s", err)
		}
	}

	return nil
}
