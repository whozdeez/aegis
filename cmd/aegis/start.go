package main

import (
	"fmt"

	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/tui"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start interactive Aegis menu",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.RunStart(); err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
