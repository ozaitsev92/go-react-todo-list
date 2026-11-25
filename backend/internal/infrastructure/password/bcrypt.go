// Package password provides the bcrypt implementation of usecase.PasswordHasher.
package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher hashes passwords with bcrypt at a fixed cost.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher returns a hasher using cost (bcrypt.MinCost..bcrypt.MaxCost).
func NewBcryptHasher(cost int) (*BcryptHasher, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return nil, fmt.Errorf("password: bcrypt cost %d out of range [%d, %d]", cost, bcrypt.MinCost, bcrypt.MaxCost)
	}

	return &BcryptHasher{cost: cost}, nil
}

// Hash returns the bcrypt hash of password.
func (h *BcryptHasher) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("password: hash: %w", err)
	}

	return string(b), nil
}

// Compare returns nil if password matches hash (constant-time).
func (h *BcryptHasher) Compare(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("password: compare: %w", err)
	}

	return nil
}
