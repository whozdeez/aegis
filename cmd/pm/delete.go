/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"

	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/input"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/security"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/storage"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/vault"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [service]",
	Short: "Delete a saved service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runDelete(args[0])
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func runDelete(service string) {
	fmt.Println()

	// ===== CHECK EXISTENCE =====
	exists, err := storage.EntryExists("data/vault.db", service)
	if err != nil {
		fmt.Println("inside err exist")
		fmt.Println("Error:", err)
		return
	}

	if !exists {
		fmt.Println("❌ Service not found")
		return
	}

	fmt.Println("Service :", service)
	fmt.Println("ℹ This action cannot be undone")

	// ===== MASTER PASSWORD =====
	master, err := input.ReadHidden("🔐 Enter master password: ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer security.ZeroBytes(master)

	// ===== VAULT META =====
	salt, checkCipher, checkNonce, err := vault.GetVaultMeta("data/vault.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// ===== INTEGRITY CHECK =====
	if err := vault.VerifyMasterPassword(master, salt, checkCipher, checkNonce); err != nil {
		fmt.Println("❌ Wrong master password")
		return
	}

	// ===== CONFIRMATION =====
	fmt.Print("\nDelete this entry? (y/N): ")
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "y" && confirm != "Y" {
		fmt.Println("ℹ Delete cancelled")
		return
	}

	// ===== DELETE =====
	if err := storage.DeleteEntry("data/vault.db", service); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("✔ Service deleted successfully")
}


