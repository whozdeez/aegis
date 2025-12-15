/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"strings"

	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/storage"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved services",
	Run: func(cmd *cobra.Command, args []string) {
		runList()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList() {
	fmt.Println()

	entries, err := storage.ListEntries("data/vault.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("ℹ No entries found")
		return
	}

	fmt.Printf("🔐 Stored entries (%d)\n\n", len(entries))

	// Header
	fmt.Printf("%-16s %s\n", "SERVICE", "USERNAME")
	fmt.Println(strings.Repeat("─", 32))

	// Rows
	for _, e := range entries {
		fmt.Printf("%-16s %s\n", e.Service, e.Username)
	}

	// Hint
	fmt.Println("\nℹ Use `pm get <service>` to view a password")
}

