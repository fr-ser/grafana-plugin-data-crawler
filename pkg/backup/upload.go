package backup

import (
	"fmt"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

// UploadDatabase uploads the provided database to Dropbox
func UploadDatabase(databaseLocation string, dropboxToken string) error {
	config := dropbox.Config{
		Token:    dropboxToken,
		LogLevel: dropbox.LogInfo,
	}
	dropboxClient := files.New(config)

	databaseFile, err := os.Open(databaseLocation)
	if err != nil {
		return fmt.Errorf("error opening database for upload: %v", err)
	}
	defer databaseFile.Close()

	_, err = dropboxClient.Upload(files.NewUploadArg("/backup/plugin.db"), databaseFile)
	if err != nil {
		return fmt.Errorf("error during file upload: %v", err)
	}
	return nil
}
