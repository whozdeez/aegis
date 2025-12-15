/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"

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
	entries, err := storage.ListEntries("data/vault.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("No entries found")
		return
	}

	fmt.Printf("Found %d entries\n\n", len(entries))
	fmt.Printf("%-12s %s\n", "SERVICE", "USERNAME")

	for _, e := range entries {
		fmt.Printf("%-12s %s\n", e.Service, e.Username)
	}
}
