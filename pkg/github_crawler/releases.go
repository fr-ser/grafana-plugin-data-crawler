package github_crawler

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/go-github/v43/github"

	// register the SQL driver
	_ "modernc.org/sqlite"
)

// StoreReleases stores the download count for releases from Github in the SQLite database
func StoreReleases(releases []*github.RepositoryRelease, databaseLocation string) error {
	db, err := sql.Open("sqlite", databaseLocation)
	if err != nil {
		return fmt.Errorf("error opending the database: %s", err)
	}
	defer db.Close()

	for _, release := range releases {
		for _, asset := range release.Assets {
			_, err = db.Exec(
				`INSERT INTO github_releases (timestamp, tag, asset_name, downloads) VALUES (?, ?, ?, ?)`,
				time.Now().Unix(), release.TagName, asset.Name, asset.DownloadCount,
			)
			if err != nil {
				return fmt.Errorf("error inserting row: %s", err)
			}
		}
	}

	return nil
}
