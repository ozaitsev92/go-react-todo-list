package task

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	// ErrTaskNotFound is returned when no task matches the query.
	ErrTaskNotFound = errors.New("task not found")
	// ErrTaskExists is returned when saving a task whose ID already exists.
	ErrTaskExists = errors.New("task already exists")
)

// Repository persists tasks. Implementations return (possibly wrapped)
// ErrTaskNotFound and ErrTaskExists; GetAllByUserID returns a non-nil slice
// ordered by CreatedAt.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (Task, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]Task, error)
	Save(ctx context.Context, t Task) error
	Update(ctx context.Context, t Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}
