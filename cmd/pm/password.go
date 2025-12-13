package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func promptMasterPassword() (string, error) {
	fmt.Print("Enter master password: ")

	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}

	return string(bytePassword), nil
}
