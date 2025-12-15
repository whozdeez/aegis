/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"

	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/crypto"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/input"
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
	// 1. Ambil data dari DB
	username, ciphertext, nonce, err := storage.GetEntry("data/vault.db", service)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 2. Prompt master password
	masterPasswordStr, err := input.ReadHidden("Enter master password: ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 3. Ambil salt
	salt, err := vault.GetVaultSalt("data/vault.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 4. Derive key
	key, err := crypto.DeriveKey(string(masterPasswordStr), salt)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 5. Decrypt password
	plaintext, err := crypto.DecryptAESGCM(key, ciphertext, nonce)
	if err != nil {
		fmt.Println("❌ Wrong master password")
		return
	}

	// 6. Tampilkan hasil
	fmt.Println("Service :", service)
	fmt.Println("Username:", username)
	fmt.Println("Password:", string(plaintext))
}
