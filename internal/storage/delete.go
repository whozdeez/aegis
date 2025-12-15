package storage

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

func DeleteEntry(dbPath, service string) error {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	result, err := db.Exec(
		`DELETE FROM entries WHERE service = ?`,
		service,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("service not found")
	}

	return nil
}
