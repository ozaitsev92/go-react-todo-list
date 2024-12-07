package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
)

var (
	ErrUnauthorizedAction = errors.New("unauthorized action")
)

type TaskUseCase struct {
	taskRepository task.Repository
}

// NewTaskUseCase creates an new instance of the TaskUseCase.
func NewTaskUseCase(taskRepository task.Repository) *TaskUseCase {
	return &TaskUseCase{
		taskRepository: taskRepository,
	}
}

// CreateTask creates a new task and saves it to the task repository.
func (s *TaskUseCase) CreateTask(ctx context.Context, text string, userId uuid.UUID) (task.Task, error) {
	t, err := task.NewTask(text, userId)
	if err != nil {
		return task.Task{}, err
	}

	err = s.taskRepository.Save(ctx, t)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// GetAllTasksForUser returns all tasks taht balong to a given user.
func (s *TaskUseCase) GetAllTasksForUser(ctx context.Context, userId uuid.UUID) ([]task.Task, error) {
	t, err := s.taskRepository.GetAllByUserID(ctx, userId)
	if err != nil {
		return []task.Task{}, err
	}

	return t, nil
}

// UpdateTask updates the task and saves it to the taskRepository.
func (s *TaskUseCase) UpdateTask(ctx context.Context, id uuid.UUID, text string, userId uuid.UUID) (task.Task, error) {
	t, err := s.taskRepository.GetByID(ctx, id)
	if err != nil {
		return task.Task{}, err
	}

	if t.UserID != userId {
		return task.Task{}, ErrUnauthorizedAction
	}

	err = t.SetText(text)
	if err != nil {
		return task.Task{}, err
	}

	err = s.taskRepository.Update(ctx, t)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// MarkTaskCompleted marks the task as completed and saves it to the taskRepository.
func (s *TaskUseCase) MarkTaskCompleted(ctx context.Context, id uuid.UUID, userId uuid.UUID) (task.Task, error) {
	t, err := s.taskRepository.GetByID(ctx, id)
	if err != nil {
		return task.Task{}, err
	}

	if t.UserID != userId {
		return task.Task{}, ErrUnauthorizedAction
	}

	t.MarkCompleted()

	err = s.taskRepository.Update(ctx, t)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// MarkTaskNotCompleted marks the task as NOT completed and saves it to the taskRepository.
func (s *TaskUseCase) MarkTaskNotCompleted(ctx context.Context, id uuid.UUID, userId uuid.UUID) (task.Task, error) {
	t, err := s.taskRepository.GetByID(ctx, id)
	if err != nil {
		return task.Task{}, err
	}

	if t.UserID != userId {
		return task.Task{}, ErrUnauthorizedAction
	}

	t.MarkNotCompleted()

	err = s.taskRepository.Update(ctx, t)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// DeleteTask deletes the task.
func (s *TaskUseCase) DeleteTask(ctx context.Context, id uuid.UUID, userId uuid.UUID) error {
	t, err := s.taskRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if t.UserID != userId {
		return ErrUnauthorizedAction
	}

	err = s.taskRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
