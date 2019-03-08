package crypto

import (
	"crypto/hmac"
	"crypto/sha256"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword salts and hashes a password using BCrypt
func HashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

// VerifyPasswordHash verifies that a password and its hash match
func VerifyPasswordHash(password string, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// HS256 performs HMAC with SHA256
func HS256(value string, key string) ([]byte, error) {
	cipher := hmac.New(sha256.New, []byte(key))
	_, err := cipher.Write([]byte(value))
	if err != nil {
		return nil, err
	}
	return cipher.Sum(nil), nil
}
