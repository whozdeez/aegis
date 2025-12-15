package vault

import (
	"errors"

	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/crypto"
)

const IntegrityCheckPlaintext = "vault-check-ok"

func VerifyMasterPassword(
	master []byte,
	salt []byte,
	checkCipher []byte,
	checkNonce []byte,
) error {

	key, err := crypto.DeriveKey(master, salt)
	if err != nil {
		return err
	}

	plaintext, err := crypto.DecryptAESGCM(key, checkCipher, checkNonce)
	if err != nil {
		return errors.New("wrong master password")
	}

	if string(plaintext) != IntegrityCheckPlaintext {
		return errors.New("vault integrity check failed")
	}

	return nil
}
