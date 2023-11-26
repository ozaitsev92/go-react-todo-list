package domain

import (
	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	TaskText  string    `json:"task_text"`
	TaskOrder int       `json:"task_order"`
	UserID    uuid.UUID `json:"user_id"`
}

type DeleteTaskRequest struct {
	ID uuid.UUID `json:"user_id"`
}

type MarkTaskDoneRequest struct {
	ID uuid.UUID `json:"user_id"`
}

type MarkTaskNotDoneRequest struct {
	ID uuid.UUID `json:"user_id"`
}

type GetTasksByUserRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type UpdateTaskRequest struct {
	ID        uuid.UUID `json:"id"`
	TaskText  string    `json:"task_text"`
	TaskOrder int       `json:"task_order"`
}

func (r *UpdateTaskRequest) EnrichTask(t *Task) {
	t.taskText = r.TaskText
	t.taskOrder = r.TaskOrder
}
