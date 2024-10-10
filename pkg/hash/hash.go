package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// Password generates a bcrypt hash for the given password.
func Password(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
