package github_crawler

import (
	"database/sql"
	"math"
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
		`CREATE TABLE github_releases (timestamp INTEGER, tag TEXT, asset_name TEXT, downloads INTEGER)`,
	)

	testReleases := []*github.RepositoryRelease{
		{
			TagName: strPointer("1.0.0"),
			Assets:  []*github.ReleaseAsset{{Name: strPointer("asset.zip"), DownloadCount: intPointer(5)}},
		},
		{
			TagName: strPointer("1.1.0"),
			Assets:  []*github.ReleaseAsset{{Name: strPointer("asset.zip"), DownloadCount: intPointer(7)}},
		},
	}

	now := time.Now().Unix()
	err := StoreReleases(testReleases, databaseLocation)

	if err != nil {
		t.Errorf("Received error: %s", err)
	}

	rows, _ := db.Query(
		"SELECT timestamp, tag, asset_name, downloads FROM github_releases ORDER BY tag ASC",
	)

	index := -1
	for rows.Next() {
		index++
		var timestamp int64
		item := github.RepositoryRelease{Assets: []*github.ReleaseAsset{{}}}
		_ = rows.Scan(&timestamp, &item.TagName, &item.Assets[0].Name, &item.Assets[0].DownloadCount)

		// as the content contains a timestamp we cannot "naively" compare for whole struct equality
		if math.Abs(float64(now-timestamp)) > 1 {
			t.Errorf("Expected timestamp of roughly %+v items. Got: %+v", now, timestamp)
		}
		if *item.TagName != *testReleases[index].TagName {
			t.Errorf("Expected tag of %+v. Got: %+v", *item.TagName, *testReleases[index].TagName)
		}
		if *item.Assets[0].Name != *testReleases[index].Assets[0].Name {
			t.Errorf("Expected asset name of %+v. Got: %+v", *item.Assets[0].Name, *testReleases[index].Assets[0].Name)
		}
		if *item.Assets[0].DownloadCount != *testReleases[index].Assets[0].DownloadCount {
			t.Errorf("Expected asset DownloadCount of %+v. Got: %+v", *item.Assets[0].DownloadCount, *testReleases[index].Assets[0].DownloadCount)
		}
	}
}

func TestGetAndStoreReleasesMultiAssets(t *testing.T) {
	tmpDir := t.TempDir()
	databaseLocation := tmpDir + "/plugin.db"

	db, _ := sql.Open("sqlite", databaseLocation)
	defer db.Close()
	_, _ = db.Exec(
		`CREATE TABLE github_releases (timestamp INTEGER, tag TEXT, asset_name TEXT, downloads INTEGER)`,
	)

	testReleases := []*github.RepositoryRelease{
		{
			TagName: strPointer("1.0.0"),
			Assets: []*github.ReleaseAsset{
				{Name: strPointer("asset1.zip"), DownloadCount: intPointer(5)},
				{Name: strPointer("asset2.zip"), DownloadCount: intPointer(7)},
			},
		},
	}

	// as the database has one row per asset, we use the below structure for the comparison
	expectedStoredReleases := []*github.RepositoryRelease{
		{
			TagName: strPointer("1.0.0"),
			Assets:  []*github.ReleaseAsset{{Name: strPointer("asset1.zip"), DownloadCount: intPointer(5)}},
		},
		{
			TagName: strPointer("1.0.0"),
			Assets:  []*github.ReleaseAsset{{Name: strPointer("asset2.zip"), DownloadCount: intPointer(7)}},
		},
	}

	now := time.Now().Unix()
	err := StoreReleases(testReleases, databaseLocation)

	if err != nil {
		t.Errorf("Received error: %s", err)
	}

	rows, _ := db.Query(
		"SELECT timestamp, tag, asset_name, downloads FROM github_releases ORDER BY asset_name ASC",
	)

	index := -1
	for rows.Next() {
		index++
		var timestamp int64
		item := github.RepositoryRelease{Assets: []*github.ReleaseAsset{{}}}
		_ = rows.Scan(&timestamp, &item.TagName, &item.Assets[0].Name, &item.Assets[0].DownloadCount)

		// as the content contains a timestamp we cannot "naively" compare for whole struct equality
		if math.Abs(float64(now-timestamp)) > 1 {
			t.Errorf("Expected timestamp of roughly %+v items. Got: %+v", now, timestamp)
		}
		if *item.TagName != *expectedStoredReleases[index].TagName {
			t.Errorf("Expected tag of %+v. Got: %+v", *item.TagName, *expectedStoredReleases[index].TagName)
		}
		if *item.Assets[0].Name != *expectedStoredReleases[index].Assets[0].Name {
			t.Errorf("Expected asset name of %+v. Got: %+v", *item.Assets[0].Name, *expectedStoredReleases[index].Assets[0].Name)
		}
		if *item.Assets[0].DownloadCount != *expectedStoredReleases[index].Assets[0].DownloadCount {
			t.Errorf("Expected asset DownloadCount of %+v. Got: %+v", *item.Assets[0].DownloadCount, *expectedStoredReleases[index].Assets[0].DownloadCount)
		}
	}
}
