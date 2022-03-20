package backup

import (
	"database/sql"
	"fmt"

	// register the SQL driver
	_ "modernc.org/sqlite"
)

// CheckDatabaseConsistency returns an error if the database is somehow malformed
func CheckDatabaseConsistency(databaseLocation string) error {
	db, err := sql.Open("sqlite", databaseLocation)
	if err != nil {
		return fmt.Errorf("error opending the database: %s", err)
	}
	defer db.Close()

	rows, err := db.Query("PRAGMA quick_check")
	if err != nil {
		return fmt.Errorf("error execting query: %s", err)
	}

	for rows.Next() {
		var result string
		err = rows.Scan(&result)
		if err != nil {
			return fmt.Errorf("error reading row: %s", err)
		}
		if result != "ok" {
			return fmt.Errorf("did not get 'ok' result but: %v", result)
		}
	}
	err = rows.Close()
	if err != nil {
		return fmt.Errorf("error closing rows: %s", err)
	}

	return nil
}
