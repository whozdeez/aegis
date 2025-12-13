package main

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

func getVaultSalt(dbPath string) ([]byte, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var salt []byte
	err = db.QueryRow("SELECT salt FROM vault_meta WHERE id = 1").Scan(&salt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("vault not initialized, run `pm init` first")
		}
		return nil, err
	}

	return salt, nil
}
