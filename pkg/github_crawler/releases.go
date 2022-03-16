package github_crawler

import (
	"database/sql"
	"fmt"

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
				`INSERT INTO github_releases (tag, asset_name, downloads, created_at) VALUES (?, ?, ?, ?)
				ON CONFLICT(tag, asset_name) DO UPDATE SET downloads=excluded.downloads`,
				release.TagName, asset.Name, asset.DownloadCount, release.CreatedAt.Unix(),
			)
			if err != nil {
				return fmt.Errorf("error inserting row: %s", err)
			}
		}
	}

	return nil
}
