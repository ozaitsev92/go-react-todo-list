package model

import (
	"time"

	"github.com/ozaitsev92/tododdd/internal/domain/task"
)

// Task -.
type Task struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponseFromTask -.
func ToResponseFromTask(t task.Task) Task {
	return Task{
		ID:        t.ID.String(),
		Text:      t.Text,
		Completed: t.Completed,
		UserID:    t.UserID.String(),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

// ToResponseFromTaskCollection -.
func ToResponseFromTaskCollection(tasks []task.Task) []Task {
	response := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		response = append(response, ToResponseFromTask(t))
	}
	return response
}
