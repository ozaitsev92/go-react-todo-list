package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/config"
	"github.com/ozaitsev92/tododdd/internal/domain/task"

	"github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/memory/converter"
	repoModel "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/memory/model"
)

var _ task.Repository = (*Repository)(nil)

type Repository struct {
	tasks map[uuid.UUID]repoModel.Task
	mu    sync.RWMutex
}

func NewRepository(_ config.Config) *Repository {
	return &Repository{
		tasks: make(map[uuid.UUID]repoModel.Task),
	}
}

func (r *Repository) GetByID(_ context.Context, id uuid.UUID) (task.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.tasks == nil {
		r.tasks = make(map[uuid.UUID]repoModel.Task)
	}

	if ti, ok := r.tasks[id]; ok {
		return converter.ToTaskFromRepo(ti), nil
	}

	return task.Task{}, task.ErrTaskNotFound
}

func (r *Repository) GetAllByUserID(_ context.Context, userId uuid.UUID) ([]task.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.tasks == nil {
		r.tasks = make(map[uuid.UUID]repoModel.Task)
	}

	var tasks []task.Task
	for _, ti := range r.tasks {
		if ti.UserID == userId.String() {
			tasks = append(tasks, converter.ToTaskFromRepo(ti))
		}
	}

	sort.Slice(tasks, func(i, j int) bool { return tasks[i].CreatedAt.UnixNano() < tasks[j].CreatedAt.UnixNano() })

	return tasks, nil
}

func (r *Repository) Save(_ context.Context, ti task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.tasks == nil {
		r.tasks = make(map[uuid.UUID]repoModel.Task)
	}

	if _, ok := r.tasks[ti.ID]; ok {
		return fmt.Errorf("task already exists: %w", task.ErrFailedToSaveTask)
	}

	r.tasks[ti.ID] = converter.ToRepoFromTask(ti)

	return nil
}

func (r *Repository) Update(_ context.Context, ti task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.tasks == nil {
		r.tasks = make(map[uuid.UUID]repoModel.Task)
	}

	if _, ok := r.tasks[ti.ID]; !ok {
		return fmt.Errorf("task does not exist: %w", task.ErrFailedToSaveTask)
	}

	r.tasks[ti.ID] = converter.ToRepoFromTask(ti)

	return nil
}

func (r *Repository) Delete(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.tasks == nil {
		r.tasks = make(map[uuid.UUID]repoModel.Task)
	}

	if _, ok := r.tasks[id]; !ok {
		return fmt.Errorf("task does not exist: %w", task.ErrFailedDeleteTask)
	}

	delete(r.tasks, id)

	return nil
}
