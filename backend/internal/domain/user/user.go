// Package user contains the user entity and its repository port.
package user

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Password length limits; bcrypt ignores bytes beyond 72.
const (
	MinPasswordLength = 8
	MaxPasswordLength = 72
)

var (
	// ErrInvalidEmail is returned for empty or malformed email addresses.
	ErrInvalidEmail = errors.New("email is invalid")
	// ErrWeakPassword is returned when a password violates the length policy.
	ErrWeakPassword = errors.New("password must be between 8 and 72 bytes")
	// ErrInvalidPasswordHash is returned when constructing a user without a hash.
	ErrInvalidPasswordHash = errors.New("password hash is required")
)

// User is the user entity.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser creates a user from an email (normalised) and a password hash.
func NewUser(email, passwordHash string) (User, error) {
	normalized, err := NormalizeEmail(email)
	if err != nil {
		return User{}, err
	}

	if passwordHash == "" {
		return User{}, ErrInvalidPasswordHash
	}

	now := now()

	return User{
		ID:           uuid.New(),
		Email:        normalized,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// NormalizeEmail trims and lower-cases a bare address, rejecting display-name
// forms such as "Alice <alice@example.com>".
func NormalizeEmail(email string) (string, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return "", ErrInvalidEmail
	}

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Name != "" || addr.Address != email {
		return "", ErrInvalidEmail
	}

	return strings.ToLower(addr.Address), nil
}

// ValidatePassword enforces the password policy.
func ValidatePassword(password string) error {
	if n := len(password); n < MinPasswordLength || n > MaxPasswordLength {
		return ErrWeakPassword
	}

	return nil
}

// now returns UTC time truncated to MongoDB's millisecond precision.
func now() time.Time {
	return time.Now().UTC().Truncate(time.Millisecond)
}
