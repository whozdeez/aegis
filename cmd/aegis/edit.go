/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"fmt"

	"github.com/GanesaAprilyanPhanama/aegis/internal/input"
	"github.com/GanesaAprilyanPhanama/aegis/internal/security"
	"github.com/GanesaAprilyanPhanama/aegis/internal/storage"
	"github.com/GanesaAprilyanPhanama/aegis/internal/vault"
	"github.com/GanesaAprilyanPhanama/aegis/internal/crypto"
	"github.com/spf13/cobra"
)


// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [service]",
	Short: "Edit username and/or password for a service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runEdit(args[0])
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func runEdit(service string) {
	// 1. Ambil entry lama
	oldUsername, ciphertext, nonce, err := storage.GetEntry("data/vault.db", service)
	if err != nil {
		fmt.Println("❌ Service not found")
		return
	}

	// 2. Ambil vault meta
	salt, checkCipher, checkNonce, err := vault.GetVaultMeta("data/vault.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 3. Input master password
	master, err := input.ReadHidden("🔐 Enter master password: ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer security.ZeroBytes(master)

	// 4. Integrity check
	if err := vault.VerifyMasterPassword(master, salt, checkCipher, checkNonce); err != nil {
		fmt.Println("❌ Wrong master password")
		return
	}

	// 5. Derive key
	key, err := crypto.DeriveKey(master, salt)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer security.ZeroBytes(key)

	// 6. Decrypt password lama
	plaintext, err := crypto.DecryptAESGCM(key, ciphertext, nonce)
	if err != nil {
		fmt.Println("❌ Vault data corrupted")
		return
	}
	defer security.ZeroBytes(plaintext)

	fmt.Println()

	// ===== USERNAME =====
	username := oldUsername
	fmt.Printf("Current username: %s\n", oldUsername)
	fmt.Print("Do you want to update the username? (y/N): ")

	var choice string
	fmt.Scanln(&choice)

	usernameChanged := false
	if choice == "y" || choice == "Y" {
		fmt.Print("New username: ")
		fmt.Scanln(&username)
		usernameChanged = username != oldUsername
	}

	fmt.Println()

	// ===== PASSWORD =====
	fmt.Println("Password is currently set.")
	fmt.Print("Do you want to update the password? (y/N): ")
	fmt.Scanln(&choice)

	password := plaintext
	passwordChanged := false
	if choice == "y" || choice == "Y" {
		newPassword, err := input.ReadHidden("New password: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer security.ZeroBytes(newPassword)

		password = newPassword
		passwordChanged = true
	}

	// ===== NO CHANGES =====
	if !usernameChanged && !passwordChanged {
		fmt.Println("\nℹ No changes made")
		return
	}

	// ===== SUMMARY =====
	fmt.Println("\nSummary of changes:")
	if usernameChanged {
		fmt.Println("- Username: updated")
	} else {
		fmt.Println("- Username: unchanged")
	}

	if passwordChanged {
		fmt.Println("- Password: updated")
	} else {
		fmt.Println("- Password: unchanged")
	}

	fmt.Print("\nApply changes? (y/N): ")
	fmt.Scanln(&choice)

	if choice != "y" && choice != "Y" {
		fmt.Println("ℹ Edit cancelled")
		return
	}

	// ===== ENCRYPT & SAVE =====
	newCipher, newNonce, err := crypto.EncryptAESGCM(key, password)
	if err != nil {
		fmt.Println("Encryption failed:", err)
		return
	}

	if err := storage.UpdateEntry(
		"data/vault.db",
		service,
		username,
		newCipher,
		newNonce,
	); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("✔ Entry updated successfully")
}
