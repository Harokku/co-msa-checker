package utils

import (
	"crypto/sha256"
	"fmt"
	"github.com/sethvargo/go-password/password"
)

// -------------------------
// Hashing functions
// -------------------------

// Hash256 Get a plain pw and return sha256 hashed string
func Hash256(plainPw string) string {
	h := sha256.New()
	h.Write([]byte(plainPw))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// NewPw generate a new random password with length digit
func NewPw(length int) (string, error) {
	res, err := password.Generate(length, 2, 0, false, false)
	if err != nil {
		return "", err
	}
	return res, nil
}
