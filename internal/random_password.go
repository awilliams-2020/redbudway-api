package internal

import (
	"crypto/rand"
	"encoding/base64"
)

// RandomPassword returns a URL-safe random string suitable for hashing as an account password.
func RandomPassword() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
