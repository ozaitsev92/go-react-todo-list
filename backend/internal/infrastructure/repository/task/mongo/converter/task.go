// Package converter maps between task.Task and its MongoDB document.
package converter

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	repoModel "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo/model"
)

// ToTaskFromRepo converts a stored document into a domain task.
func ToTaskFromRepo(t repoModel.Task) (task.Task, error) {
	id, err := uuid.Parse(t.ID)
	if err != nil {
		return task.Task{}, fmt.Errorf("tasks: invalid _id %q: %w", t.ID, err)
	}

	userID, err := uuid.Parse(t.UserID)
	if err != nil {
		return task.Task{}, fmt.Errorf("tasks: invalid user_id %q in task %s: %w", t.UserID, t.ID, err)
	}

	return task.Task{
		ID:        id,
		Text:      t.Text,
		Completed: t.Completed,
		UserID:    userID,
		CreatedAt: t.CreatedAt.UTC(),
		UpdatedAt: t.UpdatedAt.UTC(),
	}, nil
}

// ToRepoFromTask converts a domain task into its stored document.
func ToRepoFromTask(t task.Task) repoModel.Task {
	return repoModel.Task{
		ID:        t.ID.String(),
		Text:      t.Text,
		Completed: t.Completed,
		UserID:    t.UserID.String(),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
