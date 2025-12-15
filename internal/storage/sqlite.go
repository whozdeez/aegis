package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func InsertEntry(
	dbPath string,
	service string,
	username string,
	ciphertext []byte,
	nonce []byte,
) error {

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO entries (service, username, password, nonce)
		 VALUES (?, ?, ?, ?)`,
		service,
		username,
		ciphertext,
		nonce,
	)

	return err
}
