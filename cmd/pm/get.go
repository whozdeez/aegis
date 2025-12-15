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

	// Ambil data (sekali)
	username, ciphertext, nonce, err := storage.GetEntry("data/vault.db", service)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	salt, err := vault.GetVaultSalt("data/vault.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		master, err := input.ReadHidden("Enter master password: ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		defer security.ZeroBytes(master)

		key, err := crypto.DeriveKey(master, salt)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		defer security.ZeroBytes(key)

		plaintext, err := crypto.DecryptAESGCM(key, ciphertext, nonce)
		if err == nil {
			defer security.ZeroBytes(plaintext)

			fmt.Println("Service :", service)
			fmt.Println("Username:", username)
			fmt.Println("Password:", string(plaintext))
			return
		}

		fmt.Println("❌ Wrong master password")
	}

	fmt.Println("🔒 Too many failed attempts")
}


