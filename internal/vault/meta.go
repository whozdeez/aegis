package vault

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func GetVaultMeta(dbPath string) (salt, checkCipher, checkNonce []byte, err error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, nil, nil, err
	}
	defer db.Close()

	err = db.QueryRow(
		`SELECT salt, check_cipher, check_nonce FROM vault_meta WHERE id = 1`,
	).Scan(&salt, &checkCipher, &checkNonce)

	return
}
