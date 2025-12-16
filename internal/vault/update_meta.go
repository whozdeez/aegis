package vault

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

func UpdateVaultMeta(
	dbPath string,
	salt []byte,
	checkCipher []byte,
	checkNonce []byte,
) error {

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	result, err := db.Exec(`
		UPDATE vault_meta
		SET salt = ?, check_cipher = ?, check_nonce = ?
		WHERE id = 1
	`, salt, checkCipher, checkNonce)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("vault not initialized")
	}

	return nil
}
