package main

import (
	"os"

	"golang.org/x/term"
)

func readHiddenInput() ([]byte, error) {
	return term.ReadPassword(int(os.Stdin.Fd()))
}
