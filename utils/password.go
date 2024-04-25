package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(passowrd string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(passowrd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
