package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type VaultEntry struct {
	ID       int
	Service  string
	Username string
	Cipher   []byte
	Nonce    []byte
}

func GetAllEntries(dbPath string) ([]VaultEntry, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT id, service, username, password, nonce
		FROM entries
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []VaultEntry

	for rows.Next() {
		var e VaultEntry
		if err := rows.Scan(
			&e.ID,
			&e.Service,
			&e.Username,
			&e.Cipher,
			&e.Nonce,
		); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
