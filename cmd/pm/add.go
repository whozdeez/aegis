
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command

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
	// 1. Service name
	var service string
	fmt.Print("Service name: ")
	fmt.Scanln(&service)

	// 2. Username
	var username string
	fmt.Print("Username: ")
	fmt.Scanln(&username)

	// 3. Password (hidden)
	fmt.Print("Password: ")
	password, err := readHiddenInput()
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}

	// 4. Master password
	masterPasswordStr, err := promptMasterPassword()
	if err != nil {
		fmt.Println("Error reading master password:", err)
		return
	}
	masterPassword := []byte(masterPasswordStr)

	// 5. Ambil salt
	salt, err := getVaultSalt("data/vault.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 6. Derive key
	key, err := deriveKey(string(masterPassword), salt)
	if err != nil {
		fmt.Println("Key derivation failed:", err)
		return
	}

	// 7. Encrypt password
	ciphertext, nonce, err := encryptAESGCM(key, password)
	if err != nil {
		fmt.Println("Encryption failed:", err)
		return
	}

	// 8. Simpan ke DB
	err = insertEntry(
		"data/vault.db",
		service,
		username,
		ciphertext,
		nonce,
	)
	if err != nil {
		fmt.Println("Failed to save entry:", err)
		return
	}

	fmt.Println("Password saved successfully 🔐")
}