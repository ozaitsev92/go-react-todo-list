package todolist

import (
	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	TaskText  string    `json:"taskText"`
	TaskOrder int       `json:"taskOrder"`
	UserID    uuid.UUID `json:"userId"`
}

type DeleteTaskRequest struct {
	ID uuid.UUID `json:"userId"`
}

type MarkTaskDoneRequest struct {
	ID uuid.UUID `json:"userId"`
}

type MarkTaskNotDoneRequest struct {
	ID uuid.UUID `json:"userId"`
}

type GetTasksByUserRequest struct {
	UserID uuid.UUID `json:"userId"`
}

type UpdateTaskRequest struct {
	ID        uuid.UUID `json:"id"`
	TaskText  string    `json:"taskText"`
	TaskOrder int       `json:"taskOrder"`
}

func (r *UpdateTaskRequest) EnrichTask(t *Task) {
	t.taskText = r.TaskText
	t.taskOrder = r.TaskOrder
}
