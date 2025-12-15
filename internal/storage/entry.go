package storage

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

func GetEntry(dbPath, service string) (username string, ciphertext []byte, nonce []byte, err error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return "", nil, nil, err
	}
	defer db.Close()

	err = db.QueryRow(
		`SELECT username, password, nonce FROM entries WHERE service = ?`,
		service,
	).Scan(&username, &ciphertext, &nonce)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, nil, errors.New("service not found")
		}
		return "", nil, nil, err
	}

	return username, ciphertext, nonce, nil
}
