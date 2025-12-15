package storage

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

func UpdateEntry(
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

	result, err := db.Exec(
		`UPDATE entries
		 SET username = ?, password = ?, nonce = ?
		 WHERE service = ?`,
		username,
		ciphertext,
		nonce,
		service,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("service not found")
	}

	return nil
}
