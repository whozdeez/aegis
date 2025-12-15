package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func EntryExists(dbPath, service string) (bool, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow(
		`SELECT COUNT(1) FROM entries WHERE service = ?`,
		service,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
