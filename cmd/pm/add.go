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
	fmt.Println()

	// ===== SERVICE =====
	var service string
	fmt.Print("Service name: ")
	fmt.Scanln(&service)

	if service == "" {
		fmt.Println("❌ Service name cannot be empty")
		return
	}

	// ===== USERNAME =====
	var username string
	fmt.Print("Username: ")
	fmt.Scanln(&username)

	if username == "" {
		fmt.Println("❌ Username cannot be empty")
		return
	}

	// ===== PASSWORD =====
	password, err := input.ReadHidden("Password: ")
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}
	defer security.ZeroBytes(password)

	if len(password) == 0 {
		fmt.Println("❌ Password cannot be empty")
		return
	}

	fmt.Println()

	// ===== MASTER PASSWORD =====
	masterPassword, err := input.ReadHidden("🔐 Enter master password: ")
	if err != nil {
		fmt.Println("Error reading master password:", err)
		return
	}
	defer security.ZeroBytes(masterPassword)

	// ===== VAULT META =====
	salt, checkCipher, checkNonce, err := vault.GetVaultMeta("data/vault.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// ===== INTEGRITY CHECK =====
	if err := vault.VerifyMasterPassword(masterPassword, salt, checkCipher, checkNonce); err != nil {
		fmt.Println("❌ Wrong master password")
		return
	}

	// ===== SUMMARY =====
	fmt.Println("\nSummary:")
	fmt.Println("- Service :", service)
	fmt.Println("- Username:", username)
	fmt.Println("- Password: set")

	fmt.Print("\nSave this entry? (y/N): ")
	var choice string
	fmt.Scanln(&choice)

	if choice != "y" && choice != "Y" {
		fmt.Println("ℹ Add cancelled")
		return
	}

	// ===== DERIVE KEY =====
	key, err := crypto.DeriveKey(masterPassword, salt)
	if err != nil {
		fmt.Println("Key derivation failed:", err)
		return
	}
	defer security.ZeroBytes(key)

	// ===== ENCRYPT =====
	ciphertext, nonce, err := crypto.EncryptAESGCM(key, password)
	if err != nil {
		fmt.Println("Encryption failed:", err)
		return
	}

	// ===== INSERT =====
	err = storage.InsertEntry(
		"data/vault.db",
		service,
		username,
		ciphertext,
		nonce,
	)
	if err != nil {
		fmt.Println("❌ Service already exists")
		return
	}

	fmt.Println("✔ Password saved successfully 🔐")
}

