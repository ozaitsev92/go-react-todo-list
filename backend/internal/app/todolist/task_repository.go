package todolist

import (
	"context"

	"github.com/google/uuid"
)

type TaskRepository interface {
	SaveTask(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id uuid.UUID) (*Task, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*Task, error)
	DeleteTask(ctx context.Context, task *Task) error
}
