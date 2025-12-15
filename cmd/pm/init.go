/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
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

func initVault() {
	dbPath := "data/vault.db"

	if _, err := os.Stat(dbPath); err == nil {
		fmt.Println("Vault already initialized")
		return
	}

	os.MkdirAll("data", 0700)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create table
	_, err = db.Exec(`
	CREATE TABLE vault_meta (
		id INTEGER PRIMARY KEY,
		salt BLOB NOT NULL,
		check_cipher BLOB NOT NULL,
		check_nonce BLOB NOT NULL
	);
	`)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	// Generate salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}

	// Prompt master password
	master, err := input.ReadHidden("Create master password: ")
	if err != nil {
		panic(err)
	}
	defer security.ZeroBytes(master)

	// Derive key
	key, err := crypto.DeriveKey(master, salt)
	if err != nil {
		panic(err)
	}
	defer security.ZeroBytes(key)

	// Encrypt integrity check
	cipher, nonce, err := crypto.EncryptAESGCM(
		key,
		[]byte(vault.IntegrityCheckPlaintext),
	)
	if err != nil {
		panic(err)
	}

	// Save meta
	_, err = db.Exec(
		`INSERT INTO vault_meta (id, salt, check_cipher, check_nonce)
		 VALUES (1, ?, ?, ?)`,
		salt,
		cipher,
		nonce,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Vault initialized successfully 🔐")
}

