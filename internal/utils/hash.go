package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordHash), nil
}

func CompareHashAndPassword(hashedPassword, password []byte) error {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, password); err != nil {
		return err
	}

	return nil
}
