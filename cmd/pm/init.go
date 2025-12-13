/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"os"

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

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createTable := `
	CREATE TABLE vault_meta (
		id INTEGER PRIMARY KEY,
		salt BLOB NOT NULL
	);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		panic(err)
	}
	
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(
		"INSERT INTO vault_meta (id, salt) VALUES (1, ?)",
		salt,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Password vault initialized successfully")
}
