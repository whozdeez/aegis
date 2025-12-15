/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"

	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/storage"
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
	exists, err := storage.EntryExists("data/vault.db", service)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if !exists {
		fmt.Printf("Service '%s' not found\n", service)
		return
	}

	fmt.Printf("Are you sure you want to delete '%s'? (y/N): ", service)

	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "y" && confirm != "Y" {
		fmt.Println("Delete cancelled")
		return
	}

	err = storage.DeleteEntry("data/vault.db", service)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Service deleted successfully")
}

