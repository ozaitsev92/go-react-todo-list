package domain

import (
	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	TaskText  string
	TaskOrder int
	UserID    uuid.UUID
}

type DeleteTaskRequest struct {
	ID uuid.UUID
}

type MarkTaskDoneRequest struct {
	ID uuid.UUID
}

type MarkTaskNotDoneRequest struct {
	ID uuid.UUID
}

type GetTasksByUserRequest struct {
	UserID uuid.UUID
}

type UpdateTaskRequest struct {
	ID        uuid.UUID
	TaskText  string
	TaskOrder int
}

func (r *UpdateTaskRequest) EnrichTask(t *Task) {
	t.taskText = r.TaskText
	t.taskOrder = r.TaskOrder
}
