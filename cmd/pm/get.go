/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
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


var getCmd = &cobra.Command{
	Use:   "get [service]",
	Short: "Get password for a service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]
		runGet(service)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}


func runGet(service string) {
	const maxAttempts = 3

	// 1. Ambil entry (sekali)
	username, ciphertext, nonce, err := storage.GetEntry("data/vault.db", service)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 2. Ambil vault meta (sekali)
	salt, checkCipher, checkNonce, err := vault.GetVaultMeta("data/vault.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 3. Loop percobaan master password
	for attempts := 1; attempts <= maxAttempts; attempts++ {

		// 3a. Input master password
		master, err := input.ReadHidden("Enter master password: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer security.ZeroBytes(master)

		// 3b. Integrity check (VALIDASI MASTER PASSWORD)
		err = vault.VerifyMasterPassword(master, salt, checkCipher, checkNonce)
		if err != nil {
			fmt.Println("❌ Wrong master password")
			continue
		}

		// 3c. Derive key (master password SUDAH valid)
		key, err := crypto.DeriveKey(master, salt)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer security.ZeroBytes(key)

		// 3d. Decrypt password entry
		plaintext, err := crypto.DecryptAESGCM(key, ciphertext, nonce)
		if err != nil {
			fmt.Println("❌ Vault data corrupted")
			return
		}
		defer security.ZeroBytes(plaintext)

		// 3e. Tampilkan hasil
		fmt.Println("Service :", service)
		fmt.Println("Username:", username)
		fmt.Println("Password:", string(plaintext))
		return
	}

	// 4. Jika semua attempt gagal
	fmt.Println("🔒 Too many failed attempts")
}



