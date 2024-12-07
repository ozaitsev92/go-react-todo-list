package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("the user was not found in the repository")
	ErrFailedToStoreUser = errors.New("failed to store the user")
)

type Repository interface {
	GetByID(context.Context, uuid.UUID) (User, error)
	GetByEmail(context.Context, string) (User, error)
	Save(context.Context, User) error
}
