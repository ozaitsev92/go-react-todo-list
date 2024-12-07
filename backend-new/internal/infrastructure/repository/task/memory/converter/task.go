package converter

import (
	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	repoModel "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/memory/model"
)

func ToTaskFromRepo(t repoModel.Task) task.Task {
	return task.Task{
		ID:          uuid.MustParse(t.ID),
		Text:        t.Text,
		IsCompleted: t.IsCompleted,
		UserID:      uuid.MustParse(t.UserID),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func ToRepoFromTask(t task.Task) repoModel.Task {
	return repoModel.Task{
		ID:          t.ID.String(),
		Text:        t.Text,
		IsCompleted: t.IsCompleted,
		UserID:      t.UserID.String(),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
