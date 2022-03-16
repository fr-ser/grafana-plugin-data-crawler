package github_crawler

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/v43/github"
)

func strPointer(x string) *string {
	return &x
}

func TestGetAndStoreReleasesMultiReleases(t *testing.T) {
	tmpDir := t.TempDir()
	databaseLocation := tmpDir + "/plugin.db"

	db, _ := sql.Open("sqlite", databaseLocation)
	defer db.Close()
	_, _ = db.Exec(
		`CREATE TABLE github_releases (
			tag TEXT, asset_name TEXT, downloads INTEGER, created_at INTEGER, UNIQUE (tag, asset_name)
		)`,
	)

	testReleases := []*github.RepositoryRelease{
		{
			TagName:   strPointer("1.0.0"),
			CreatedAt: &github.Timestamp{Time: time.Unix(123, 0)},
			Assets:    []*github.ReleaseAsset{{Name: strPointer("asset.zip"), DownloadCount: intPointer(5)}},
		},
		{
			TagName:   strPointer("1.1.0"),
			CreatedAt: &github.Timestamp{Time: time.Unix(456, 0)},
			Assets:    []*github.ReleaseAsset{{Name: strPointer("asset.zip"), DownloadCount: intPointer(7)}},
		},
	}

	err := StoreReleases(testReleases, databaseLocation)

	if err != nil {
		t.Errorf("Received error: %s", err)
	}

	rows, _ := db.Query(
		"SELECT tag, asset_name, downloads, created_at FROM github_releases ORDER BY tag ASC",
	)
	storedReleases := []*github.RepositoryRelease{}

	for rows.Next() {
		var timestamp int64
		item := github.RepositoryRelease{Assets: []*github.ReleaseAsset{{}}}
		_ = rows.Scan(&item.TagName, &item.Assets[0].Name, &item.Assets[0].DownloadCount, &timestamp)

		item.CreatedAt = &github.Timestamp{Time: time.Unix(timestamp, 0)}

		storedReleases = append(storedReleases, &item)
	}

	if !reflect.DeepEqual(storedReleases, testReleases) {
		t.Errorf("Expected %+v items. Got: %+v", testReleases, storedReleases)
	}
}

func TestGetAndStoreReleasesMultiAssets(t *testing.T) {
	tmpDir := t.TempDir()
	databaseLocation := tmpDir + "/plugin.db"

	db, _ := sql.Open("sqlite", databaseLocation)
	defer db.Close()
	_, _ = db.Exec(
		`CREATE TABLE github_releases (
			tag TEXT, asset_name TEXT, downloads INTEGER, created_at INTEGER, UNIQUE (tag, asset_name)
		)`,
	)

	testReleases := []*github.RepositoryRelease{
		{
			TagName:   strPointer("1.0.0"),
			CreatedAt: &github.Timestamp{Time: time.Unix(123, 0)},
			Assets: []*github.ReleaseAsset{
				{Name: strPointer("asset1.zip"), DownloadCount: intPointer(5)},
				{Name: strPointer("asset2.zip"), DownloadCount: intPointer(7)},
			},
		},
	}

	expectedStoredReleases := []*github.RepositoryRelease{
		{
			TagName:   strPointer("1.0.0"),
			CreatedAt: &github.Timestamp{Time: time.Unix(123, 0)},
			Assets:    []*github.ReleaseAsset{{Name: strPointer("asset1.zip"), DownloadCount: intPointer(5)}},
		},
		{
			TagName:   strPointer("1.0.0"),
			CreatedAt: &github.Timestamp{Time: time.Unix(123, 0)},
			Assets:    []*github.ReleaseAsset{{Name: strPointer("asset2.zip"), DownloadCount: intPointer(7)}},
		},
	}

	err := StoreReleases(testReleases, databaseLocation)

	if err != nil {
		t.Errorf("Received error: %s", err)
	}

	rows, _ := db.Query(
		"SELECT tag, asset_name, downloads, created_at FROM github_releases ORDER BY asset_name ASC",
	)
	storedReleases := []*github.RepositoryRelease{}

	for rows.Next() {
		var timestamp int64
		item := github.RepositoryRelease{Assets: []*github.ReleaseAsset{{}}}
		_ = rows.Scan(&item.TagName, &item.Assets[0].Name, &item.Assets[0].DownloadCount, &timestamp)

		item.CreatedAt = &github.Timestamp{Time: time.Unix(timestamp, 0)}

		storedReleases = append(storedReleases, &item)
	}

	if !reflect.DeepEqual(storedReleases, expectedStoredReleases) {
		t.Errorf("Expected %+v items. Got: %+v", expectedStoredReleases, storedReleases)
	}
}

func TestGetAndStoreReleasesUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	databaseLocation := tmpDir + "/plugin.db"

	db, _ := sql.Open("sqlite", databaseLocation)
	defer db.Close()
	_, _ = db.Exec(
		`CREATE TABLE github_releases (
			tag TEXT, asset_name TEXT, downloads INTEGER, created_at INTEGER, UNIQUE (tag, asset_name)
		)`,
	)
	_, _ = db.Exec(
		"INSERT INTO github_releases (tag, asset_name, downloads, created_at) VALUES ('1.0.0', 'asset.zip', 1, 123)",
	)

	testReleases := []*github.RepositoryRelease{
		{
			TagName:   strPointer("1.0.0"),
			CreatedAt: &github.Timestamp{Time: time.Unix(123, 0)},
			Assets:    []*github.ReleaseAsset{{Name: strPointer("asset.zip"), DownloadCount: intPointer(5)}},
		},
		{
			TagName:   strPointer("1.1.0"),
			CreatedAt: &github.Timestamp{Time: time.Unix(456, 0)},
			Assets:    []*github.ReleaseAsset{{Name: strPointer("asset.zip"), DownloadCount: intPointer(7)}},
		},
	}
	err := StoreReleases(testReleases, databaseLocation)

	if err != nil {
		t.Errorf("Received error: %s", err)
	}

	rows, _ := db.Query(
		"SELECT tag, asset_name, downloads, created_at FROM github_releases ORDER BY tag ASC",
	)
	storedReleases := []*github.RepositoryRelease{}

	for rows.Next() {
		var timestamp int64
		item := github.RepositoryRelease{Assets: []*github.ReleaseAsset{{}}}
		_ = rows.Scan(&item.TagName, &item.Assets[0].Name, &item.Assets[0].DownloadCount, &timestamp)

		item.CreatedAt = &github.Timestamp{Time: time.Unix(timestamp, 0)}

		storedReleases = append(storedReleases, &item)
	}

	if !reflect.DeepEqual(storedReleases, testReleases) {
		t.Errorf("Expected %+v items. Got: %+v", testReleases, storedReleases)
	}
}
