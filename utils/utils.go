package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	length   = uint32(len(alphabet))
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateAPIKey gives an API key for a user
func GenerateAPIKey() string {
	b, err := generateRandomBytes(16)
	for err != nil {
		b, err = generateRandomBytes(16)
	}
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// GenerateFileName gives random file name restricted to 4 bytes
func GenerateFileName() string {
	seed, err := generateRandomBytes(4)
	for err != nil {
		seed, err = generateRandomBytes(4)
	}
	n := binary.BigEndian.Uint32(seed)
	if n == 0 {
		return "0"
	}

	b := make([]byte, 0, 512)
	for n > 0 {
		r := n % length
		n /= length
		b = append([]byte{alphabet[r]}, b...)
	}
	return string(b)
}
