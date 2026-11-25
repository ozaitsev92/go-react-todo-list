package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	// ErrUserNotFound is returned when no user matches the query.
	ErrUserNotFound = errors.New("user not found")
	// ErrEmailTaken is returned when saving a user whose email already exists.
	ErrEmailTaken = errors.New("email is already registered")
)

// Repository persists users. Implementations return (possibly wrapped)
// ErrUserNotFound and ErrEmailTaken; emails are already normalised.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Save(ctx context.Context, u User) error
}
