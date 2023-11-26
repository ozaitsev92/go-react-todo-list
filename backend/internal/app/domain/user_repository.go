package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}
