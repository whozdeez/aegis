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

	fmt.Println()

	// ===== GET ENTRY =====
	username, ciphertext, nonce, err := storage.GetEntry("data/vault.db", service)
	if err != nil {
		fmt.Println("❌ Service not found")
		return
	}

	// ===== VAULT META =====
	salt, checkCipher, checkNonce, err := vault.GetVaultMeta("data/vault.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// ===== AUTH LOOP =====
	for attempts := 1; attempts <= maxAttempts; attempts++ {

		master, err := input.ReadHidden("🔐 Enter master password: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// --- integrity check ---
		if err := vault.VerifyMasterPassword(master, salt, checkCipher, checkNonce); err != nil {
			fmt.Println("❌ Wrong master password")
			security.ZeroBytes(master)
			continue
		}

		// --- derive key ---
		key, err := crypto.DeriveKey(master, salt)
		security.ZeroBytes(master)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// --- decrypt ---
		plaintext, err := crypto.DecryptAESGCM(key, ciphertext, nonce)
		security.ZeroBytes(key)
		if err != nil {
			fmt.Println("❌ Vault data corrupted")
			return
		}

		// ===== SUCCESS =====
		fmt.Println("\n✔ Access granted")
		fmt.Println("Service :", service)
		fmt.Println("Username:", username)

		// ===== SHOW PASSWORD CONFIRMATION =====
		fmt.Print("\nShow password? (y/N): ")
		var choice string
		fmt.Scanln(&choice)

		if choice == "y" || choice == "Y" {
			fmt.Println("Password:", string(plaintext))
		} else {
			fmt.Println("ℹ Password hidden")
		}

		security.ZeroBytes(plaintext)
		return
	}

	// ===== LOCKOUT =====
	fmt.Println("🔒 Too many failed attempts")
}





