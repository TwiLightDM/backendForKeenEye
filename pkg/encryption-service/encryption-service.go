package encryption_service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type EncryptionService struct {
	SaltLength int
}

func NewEncryptionService(salt int) *EncryptionService {
	return &EncryptionService{
		SaltLength: salt,
	}
}

func (e EncryptionService) HashPassword(password string) (string, string, error) {
	salt, err := e.saltGeneration()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate salt: %w", err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(salt+password), bcrypt.MinCost)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), salt, nil
}

func (e EncryptionService) PasswordComparison(hashedPassword, password, salt string) (bool, error) {
	saltPassword := salt + password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltPassword))
	if err != nil {
		return false, fmt.Errorf("invalid password: %w", err)
	}
	return true, nil
}

func (e EncryptionService) saltGeneration() (string, error) {
	bytes := make([]byte, e.SaltLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(bytes)[:e.SaltLength], nil
}
