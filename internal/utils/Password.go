package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("Cant generate hash function %w", err)
	}
	return string(hashedPass), nil

}
func VerifyPassword(hashedPass string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	if err != nil {
		return fmt.Errorf("Verifying password %w", err)
	}
	return nil
}
