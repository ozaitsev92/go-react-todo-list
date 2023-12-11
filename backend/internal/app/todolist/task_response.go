package todolist

import (
	"time"

	"github.com/google/uuid"
)

type TaskResponse struct {
	ID        uuid.UUID `json:"id"`
	TaskText  string    `json:"taskText"`
	TaskOrder int       `json:"order"`
	IsDone    bool      `json:"isDone"`
	UserID    uuid.UUID `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
