// Package auth handles password hashing and credential checks.
package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost     = 12
	MinPasswordLen = 8
	MaxPasswordLen = 64
)

var ErrPasswordLength = errors.New(
	"password must be between 8 and 64 characters",
)

// HashPassword creates a bcrypt hash suitable for storage.
func HashPassword(plain string) (string, error) {
	if len(plain) < MinPasswordLen || len(plain) > MaxPasswordLen {
		return "", ErrPasswordLength
	}

	b, err := bcrypt.GenerateFromPassword(
		[]byte(plain),
		bcryptCost,
	)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// CheckPassword checks whether the password matches the stored hash.
func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(plain),
	) == nil
}
