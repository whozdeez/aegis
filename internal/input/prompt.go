package input

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func ReadHidden(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	pw, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return pw, err
	}
	fmt.Println()
	return pw, err
}
