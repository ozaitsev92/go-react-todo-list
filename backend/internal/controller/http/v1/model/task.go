// Package model contains the JSON DTOs of the v1 API.
package model

import (
	"time"

	"github.com/ozaitsev92/tododdd/internal/domain/task"
)

// Task is the JSON representation of a task.
type Task struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponseFromTask converts a domain task into its JSON representation.
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

// ToResponseFromTaskCollection converts tasks, always returning a non-nil slice.
func ToResponseFromTaskCollection(tasks []task.Task) []Task {
	response := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		response = append(response, ToResponseFromTask(t))
	}
	return response
}
