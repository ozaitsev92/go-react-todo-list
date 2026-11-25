// Package usecase contains the application services.
package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
)

// ErrUnauthorizedAction is returned when a user acts on another user's task.
var ErrUnauthorizedAction = errors.New("unauthorized action")

// TaskUseCase manages a user's tasks.
type TaskUseCase struct {
	tasks task.Repository
}

// NewTaskUseCase creates a TaskUseCase backed by tasks.
func NewTaskUseCase(tasks task.Repository) *TaskUseCase {
	return &TaskUseCase{tasks: tasks}
}

// CreateTask creates a new task owned by userID.
func (s *TaskUseCase) CreateTask(ctx context.Context, text string, userID uuid.UUID) (task.Task, error) {
	t, err := task.NewTask(text, userID)
	if err != nil {
		return task.Task{}, err
	}

	if err := s.tasks.Save(ctx, t); err != nil {
		return task.Task{}, fmt.Errorf("create task: %w", err)
	}

	return t, nil
}

// GetAllTasksForUser returns the user's tasks, oldest first.
func (s *TaskUseCase) GetAllTasksForUser(ctx context.Context, userID uuid.UUID) ([]task.Task, error) {
	tasks, err := s.tasks.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list tasks for user %s: %w", userID, err)
	}

	return tasks, nil
}

// GetTask returns a task owned by userID.
func (s *TaskUseCase) GetTask(ctx context.Context, id uuid.UUID, userID uuid.UUID) (task.Task, error) {
	return s.ownedTask(ctx, id, userID)
}

// UpdateTask changes the text of a task owned by userID.
func (s *TaskUseCase) UpdateTask(ctx context.Context, id uuid.UUID, text string, userID uuid.UUID) (task.Task, error) {
	return s.mutate(ctx, id, userID, func(t *task.Task) error { return t.SetText(text) })
}

// SetTaskCompleted marks a task owned by userID as completed or not completed.
func (s *TaskUseCase) SetTaskCompleted(ctx context.Context, id uuid.UUID, completed bool, userID uuid.UUID) (task.Task, error) {
	return s.mutate(ctx, id, userID, func(t *task.Task) error {
		t.SetCompleted(completed)
		return nil
	})
}

// DeleteTask deletes a task owned by userID.
func (s *TaskUseCase) DeleteTask(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if _, err := s.ownedTask(ctx, id, userID); err != nil {
		return err
	}

	if err := s.tasks.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete task %s: %w", id, err)
	}

	return nil
}

func (s *TaskUseCase) mutate(ctx context.Context, id, userID uuid.UUID, change func(*task.Task) error) (task.Task, error) {
	t, err := s.ownedTask(ctx, id, userID)
	if err != nil {
		return task.Task{}, err
	}

	if err := change(&t); err != nil {
		return task.Task{}, err
	}

	if err := s.tasks.Update(ctx, t); err != nil {
		return task.Task{}, fmt.Errorf("update task %s: %w", id, err)
	}

	return t, nil
}

func (s *TaskUseCase) ownedTask(ctx context.Context, id, userID uuid.UUID) (task.Task, error) {
	t, err := s.tasks.GetByID(ctx, id)
	if err != nil {
		return task.Task{}, fmt.Errorf("get task %s: %w", id, err)
	}

	if !t.IsOwnedBy(userID) {
		return task.Task{}, ErrUnauthorizedAction
	}

	return t, nil
}
