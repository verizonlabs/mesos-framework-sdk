package utils

import (
	"crypto/rand"
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

func UuidAsString() string {
	uuid := Uuid()
	return fmt.Sprintf("%X-%X-%X-%X-%X", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
