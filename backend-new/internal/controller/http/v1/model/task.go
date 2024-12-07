package model

import (
	"time"

	"github.com/ozaitsev92/tododdd/internal/domain/task"
)

type Task struct {
	ID          string    `json:"id"`
	Text        string    `json:"text"`
	IsCompleted bool      `json:"is_completed"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToResponseFromTask(t task.Task) Task {
	return Task{
		ID:          t.ID.String(),
		Text:        t.Text,
		IsCompleted: t.IsCompleted,
		UserID:      t.UserID.String(),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
