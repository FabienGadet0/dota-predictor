package helpers

import (
	"crypto/rand"
	"fmt"
)

// TokenGenerator for user
func TokenGenerator() (string, error) {
	r := make([]byte, 16)

	_, err := rand.Read(r)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", r), nil
}
