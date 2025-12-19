/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"bytes"
	"crypto/rand"
	"fmt"

	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/crypto"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/input"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/security"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/storage"
	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/vault"
	"github.com/spf13/cobra"
)

var changeMasterCmd = &cobra.Command{
	Use:   "change-master",
	Short: "Change the master password of the vault",
	Run: func(cmd *cobra.Command, args []string) {
		runChangeMaster()
	},
}


func init() {
	rootCmd.AddCommand(changeMasterCmd)
}

func runChangeMaster() {
	fmt.Println()

	// ===== VAULT META (OLD) =====
	oldSalt, checkCipher, checkNonce, err := vault.GetVaultMeta("data/vault.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// ===== MASTER PASSWORD LAMA =====
	oldMaster, err := input.ReadHidden("🔐 Enter current master password: ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer security.ZeroBytes(oldMaster)

	// ===== VERIFY OLD MASTER =====
	if err := vault.VerifyMasterPassword(oldMaster, oldSalt, checkCipher, checkNonce); err != nil {
		fmt.Println("❌ Wrong master password")
		return
	}

	// ===== DERIVE OLD KEY =====
	oldKey, err := crypto.DeriveKey(oldMaster, oldSalt)
	if err != nil {
		fmt.Println("Key derivation failed:", err)
		return
	}
	defer security.ZeroBytes(oldKey)

	// ===== LOAD ALL ENTRIES =====
	entries, err := storage.GetAllEntries("data/vault.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// ===== DECRYPT ALL ENTRIES =====
	type decryptedEntry struct {
		ID       int
		Password []byte
	}

	var decrypted []decryptedEntry

	for _, e := range entries {
		plaintext, err := crypto.DecryptAESGCM(oldKey, e.Cipher, e.Nonce)
		if err != nil {
			fmt.Println("❌ Vault data corrupted")
			return
		}

		decrypted = append(decrypted, decryptedEntry{
			ID:       e.ID,
			Password: plaintext,
		})
	}

	// ===== MASTER PASSWORD BARU =====
	fmt.Println()
	newMaster, err := input.ReadHidden("🔐 Enter new master password: ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer security.ZeroBytes(newMaster)

	confirm, err := input.ReadHidden("🔐 Confirm new master password: ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer security.ZeroBytes(confirm)

	if !bytes.Equal(newMaster, confirm) {
		fmt.Println("❌ Master passwords do not match")
		return
	}

	// ===== CONFIRMATION =====
	fmt.Println("\n⚠️ This will re-encrypt all stored passwords.")
	fmt.Print("Continue? (y/N): ")
	var choice string
	fmt.Scanln(&choice)

	if choice != "y" && choice != "Y" {
		fmt.Println("ℹ Operation cancelled")
		return
	}

	// ===== GENERATE NEW SALT =====
	newSalt := make([]byte, 16)
	if _, err := rand.Read(newSalt); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// ===== DERIVE NEW KEY =====
	newKey, err := crypto.DeriveKey(newMaster, newSalt)
	if err != nil {
		fmt.Println("Key derivation failed:", err)
		return
	}
	defer security.ZeroBytes(newKey)

	// ===== RE-ENCRYPT INTEGRITY CHECK =====
	newCheckCipher, newCheckNonce, err := crypto.EncryptAESGCM(
		newKey,
		[]byte(vault.IntegrityCheckPlaintext),
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// ===== RE-ENCRYPT ALL ENTRIES =====
	for i, _ := range entries {
		plain := decrypted[i].Password

		cipher, nonce, err := crypto.EncryptAESGCM(newKey, plain)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// overwrite entry cipher
		entries[i].Cipher = cipher
		entries[i].Nonce = nonce

		security.ZeroBytes(plain)
	}

	// ===== UPDATE VAULT META =====
	if err := vault.UpdateVaultMeta(
		"data/vault.db",
		newSalt,
		newCheckCipher,
		newCheckNonce,
	); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// ===== UPDATE ALL ENTRIES (ATOMIC) =====
	if err := storage.UpdateAllEntries("data/vault.db", entries); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("✔ Master password changed successfully 🔐")
}

