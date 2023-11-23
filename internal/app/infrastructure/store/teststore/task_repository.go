package teststore

import (
	"context"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/domain"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store"
)

type TaskRepository struct {
	store *Store
	tasks map[uuid.UUID]*domain.Task
}

func (r *TaskRepository) SaveTask(ctx context.Context, task *domain.Task) error {
	r.tasks[task.GetID()] = task
	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	if t, ok := r.tasks[id]; ok {
		return t, nil
	} else {
		return nil, store.ErrRecordNotFound
	}
}

func (r *TaskRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Task, error) {
	tasks := []*domain.Task{}
	for _, t := range r.tasks {
		if t.GetUserID() == userID {
			tasks = append(tasks, t)
		}
	}

	return tasks, nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, task *domain.Task) error {
	if _, ok := r.tasks[task.GetID()]; ok {
		delete(r.tasks, task.GetID())
		return nil
	} else {
		return store.ErrRecordNotFound
	}
}
