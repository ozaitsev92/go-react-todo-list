package task

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTaskNotFound     = errors.New("the task was not found in the repository")
	ErrFailedToSaveTask = errors.New("failed to save the task")
	ErrFailedUpdateTask = errors.New("failed to update the task")
	ErrFailedDeleteTask = errors.New("failed to delete the task")
)

type Repository interface {
	GetByID(context.Context, uuid.UUID) (Task, error)
	GetAllByUserID(context.Context, uuid.UUID) ([]Task, error)
	Save(context.Context, Task) error
	Update(context.Context, Task) error
	Delete(context.Context, uuid.UUID) error
}
