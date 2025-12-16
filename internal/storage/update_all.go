package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func UpdateAllEntries(dbPath string, entries []VaultEntry) error {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		UPDATE entries
		SET password = ?, nonce = ?
		WHERE id = ?
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, e := range entries {
		if _, err := stmt.Exec(
			e.Cipher,
			e.Nonce,
			e.ID,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
