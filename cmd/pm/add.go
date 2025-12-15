package main

import (
	"fmt"

	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/crypto"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/input"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/security"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/storage"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/vault"
	"github.com/spf13/cobra"
)


var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new password entry",
	Run: func(cmd *cobra.Command, args []string) {
		runAdd()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd() {
	// 1. Service name
	var service string
	fmt.Print("Service name: ")
	fmt.Scanln(&service)

	// 2. Username
	var username string
	fmt.Print("Username: ")
	fmt.Scanln(&username)

	// 3. Password (hidden)
	password, err := input.ReadHidden("Password: ")
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}
	defer security.ZeroBytes(password)

	// 4. Master password
	masterPassword, err := input.ReadHidden("Enter master password: ")
	if err != nil {
		fmt.Println("Error reading master password:", err)
		return
	}
	defer security.ZeroBytes(masterPassword)

	// 5. Ambil salt
	salt, err := vault.GetVaultSalt("data/vault.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 6. Derive key (TANPA string)
	key, err := crypto.DeriveKey(masterPassword, salt)
	if err != nil {
		fmt.Println("Key derivation failed:", err)
		return
	}
	defer security.ZeroBytes(key)

	// 7. Encrypt password
	ciphertext, nonce, err := crypto.EncryptAESGCM(key, password)
	if err != nil {
		fmt.Println("Encryption failed:", err)
		return
	}

	// 8. Simpan ke DB
	err = storage.InsertEntry(
		"data/vault.db",
		service,
		username,
		ciphertext,
		nonce,
	)
	if err != nil {
		fmt.Println("Failed to save entry:", err)
		return
	}

	fmt.Println("Password saved successfully 🔐")
}
