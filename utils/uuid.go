package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
)

// Generates a UUID using random bytes from a secure source.
func Uuid() []byte {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic("Failed to generate UUID")
	}

	// Generate a v4 UUID.
	b[6] = (b[6] | 0x40) & 0x4F
	b[8] = (b[8] | 0x80) & 0xBF

	return b
}

// Converts a UUID in bytes to a string.
func UuidToString(uuid []byte) (string, error) {
	var err error
	var id string
	if len(uuid) > 0 && len(uuid) <= 16 {
		id = fmt.Sprintf("%X-%X-%X-%X-%X", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
	} else {
		err = errors.New("Invalid UUID passed into UuidToString")
	}
	return id, err
}

func UuidAsString() (string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic("Failed to generate UUID")
	}

	// Generate a v4 UUID.
	b[6] = (b[6] | 0x40) & 0x4F
	b[8] = (b[8] | 0x80) & 0xBF

	uuid := b

	return fmt.Sprintf("%X-%X-%X-%X-%X", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}