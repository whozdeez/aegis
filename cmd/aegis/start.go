package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/GanesaAprilyanPhanama/passmanager-cli/internal/tui"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start interactive Aegis menu",
	Run: func(cmd *cobra.Command, args []string) {

		for {
			clearScreen() // ① sebelum TUI

			action, err := tui.RunStart()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			clearScreen() // ② setelah TUI, sebelum action

			switch action {

			case tui.ActionAdd:
				runAdd()
				waitForEnter()

			case tui.ActionGet:
				var service string
				fmt.Print("Service name: ")
				fmt.Scanln(&service)
				runGet(service)
				waitForEnter()

			case tui.ActionEdit:
				var service string
				fmt.Print("Service name: ")
				fmt.Scanln(&service)
				runEdit(service)
				waitForEnter()

			case tui.ActionDelete:
				var service string
				fmt.Print("Service name: ")
				fmt.Scanln(&service)
				runDelete(service)
				waitForEnter()

			case tui.ActionList:
				runList()
				waitForEnter()

			case tui.ActionChangeMaster:
				runChangeMaster()
				waitForEnter()

			case tui.ActionExit:
				clearScreen() // ③ sebelum exit
				fmt.Println("Goodbye 👋")
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func waitForEnter() {
	fmt.Println("\nPress Enter to return to menu...")
	fmt.Scanln()
}

func clearScreen() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default: 
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
