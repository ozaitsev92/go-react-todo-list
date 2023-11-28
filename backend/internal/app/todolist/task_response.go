package todolist

import (
	"time"

	"github.com/google/uuid"
)

type TaskResponse struct {
	ID        uuid.UUID `json:"id"`
	TaskText  string    `json:"task_text"`
	TaskOrder int       `json:"order"`
	IsDone    bool      `json:"is_done"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func TaskToResponse(t Task) TaskResponse {
	return TaskResponse{
		ID:        t.id,
		TaskText:  t.taskText,
		TaskOrder: t.taskOrder,
		IsDone:    t.isDone,
		UserID:    t.userID,
		CreatedAt: t.createdAt,
		UpdatedAt: t.updatedAt,
	}
}
