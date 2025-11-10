package utils

import "golang.org/x/crypto/bcrypt"

// Parolni xeshlovchi funksiya
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
