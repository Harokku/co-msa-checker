package utils

import (
	"fmt"
	"os"
)

//ReadEnv read s env variable and return it or error if failed
func ReadEnv(s string) (string, error) {
	r := os.Getenv(s)
	if r == "" {
		return "", fmt.Errorf("passed env variable (%s) not set", s)
	}
	return r, nil
}
