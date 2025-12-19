/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"fmt"
	"os"

	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/crypto"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/input"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/security"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/vault"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize password vault",
	Run: func(cmd *cobra.Command, args []string) {
		initVault()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func promptValidMasterPassword() ([]byte, error) {
	for {
		master, err := input.ReadHidden("🔐 Create master password: ")
		if err != nil {
			return nil, err
		}

		if len(master) < 8 {
			fmt.Println("❌ Master password must be at least 8 characters")
			security.ZeroBytes(master)
			continue
		}

		confirm, err := input.ReadHidden("🔐 Confirm master password: ")
		if err != nil {
			security.ZeroBytes(master)
			return nil, err
		}

		if !bytes.Equal(master, confirm) {
			fmt.Println("❌ Master passwords do not match")
			security.ZeroBytes(master)
			security.ZeroBytes(confirm)
			continue
		}

		security.ZeroBytes(confirm)
		return master, nil
	}
}


func initVault() {
	dbPath := "data/vault.db"

	if _, err := os.Stat(dbPath); err == nil {
		fmt.Println("ℹ Vault already initialized")
		return
	}

	fmt.Println()
	fmt.Println("⚠️ IMPORTANT NOTICE")
	fmt.Println("This password vault uses zero-knowledge encryption.")
	fmt.Println()
	fmt.Println("If you forget your master password:")
	fmt.Println("- Your passwords CANNOT be recovered")
	fmt.Println("- There is NO reset or backdoor")
	fmt.Println("- The vault will be permanently inaccessible")
	fmt.Println()

	fmt.Print("Continue initialization? (y/N): ")
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "y" && confirm != "Y" {
		fmt.Println("ℹ Initialization cancelled")
		return
	}

	// ===== MASTER PASSWORD (LOOP UNTIL VALID) =====
	master, err := promptValidMasterPassword()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer security.ZeroBytes(master)

	// ===== FILESYSTEM =====
	if err := os.MkdirAll("data", 0700); err != nil {
		fmt.Println("Error:", err)
		return
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer db.Close()

	// ===== CREATE TABLES =====
	_, err = db.Exec(`
	CREATE TABLE vault_meta (
		id INTEGER PRIMARY KEY,
		salt BLOB NOT NULL,
		check_cipher BLOB NOT NULL,
		check_nonce BLOB NOT NULL
	);
	`)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_, err = db.Exec(`
	CREATE TABLE entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		service TEXT UNIQUE NOT NULL,
		username TEXT NOT NULL,
		password BLOB NOT NULL,
		nonce BLOB NOT NULL
	);
	`)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// ===== SALT =====
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// ===== DERIVE KEY =====
	key, err := crypto.DeriveKey(master, salt)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer security.ZeroBytes(key)

	// ===== INTEGRITY CHECK =====
	cipher, nonce, err := crypto.EncryptAESGCM(
		key,
		[]byte(vault.IntegrityCheckPlaintext),
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// ===== SAVE META =====
	_, err = db.Exec(
		`INSERT INTO vault_meta (id, salt, check_cipher, check_nonce)
		 VALUES (1, ?, ?, ?)`,
		salt,
		cipher,
		nonce,
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("\n✔ Vault initialized successfully 🔐")
}



