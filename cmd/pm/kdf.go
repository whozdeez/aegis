package main

import (
	"golang.org/x/crypto/scrypt"
)

func deriveKey(masterPassword string, salt []byte) ([]byte, error) {
	const (
		N = 1 << 15 // CPU/memory cost
		r = 8
		p = 1
		keyLen = 32 // AES-256
	)

	return scrypt.Key(
		[]byte(masterPassword),
		salt,
		N,
		r,
		p,
		keyLen,
	)
}
