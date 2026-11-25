// Package repository is a concurrency-safe in-memory task.Repository.
package repository

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
)

var _ task.Repository = (*Repository)(nil)

// Repository stores tasks by ID; tasks are values, so callers get copies.
type Repository struct {
	mu    sync.RWMutex
	tasks map[uuid.UUID]task.Task
}

// NewRepository returns an empty repository.
func NewRepository() *Repository {
	return &Repository{tasks: make(map[uuid.UUID]task.Task)}
}

// GetByID returns the task with id or task.ErrTaskNotFound.
func (r *Repository) GetByID(_ context.Context, id uuid.UUID) (task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tasks[id]
	if !ok {
		return task.Task{}, task.ErrTaskNotFound
	}

	return t, nil
}

// GetAllByUserID returns the user's tasks ordered by CreatedAt (never nil).
func (r *Repository) GetAllByUserID(_ context.Context, userID uuid.UUID) ([]task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]task.Task, 0)
	for _, t := range r.tasks {
		if t.UserID == userID {
			tasks = append(tasks, t)
		}
	}

	slices.SortStableFunc(tasks, func(a, b task.Task) int { return a.CreatedAt.Compare(b.CreatedAt) })

	return tasks, nil
}

// Save inserts a new task or returns task.ErrTaskExists.
func (r *Repository) Save(_ context.Context, t task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[t.ID]; ok {
		return fmt.Errorf("save task %s: %w", t.ID, task.ErrTaskExists)
	}

	r.tasks[t.ID] = t

	return nil
}

// Update replaces an existing task or returns task.ErrTaskNotFound.
func (r *Repository) Update(_ context.Context, t task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[t.ID]; !ok {
		return task.ErrTaskNotFound
	}

	r.tasks[t.ID] = t

	return nil
}

// Delete removes a task or returns task.ErrTaskNotFound.
func (r *Repository) Delete(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[id]; !ok {
		return task.ErrTaskNotFound
	}

	delete(r.tasks, id)

	return nil
}
