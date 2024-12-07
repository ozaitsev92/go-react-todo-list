package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/config"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/memory/converter"
	repoModel "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/memory/model"
)

var _ user.Repository = (*Repository)(nil)

type Repository struct {
	users map[uuid.UUID]repoModel.User
	mu    sync.RWMutex
}

func NewRepository(_ config.Config) *Repository {
	return &Repository{
		users: make(map[uuid.UUID]repoModel.User),
	}
}

func (r *Repository) GetByID(_ context.Context, id uuid.UUID) (user.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.users == nil {
		r.users = make(map[uuid.UUID]repoModel.User)
	}

	if u, ok := r.users[id]; ok {
		return converter.ToUserFromRepo(u), nil
	}

	return user.User{}, user.ErrUserNotFound
}

func (r *Repository) GetByEmail(_ context.Context, email string) (user.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.users == nil {
		r.users = make(map[uuid.UUID]repoModel.User)
	}

	for _, u := range r.users {
		if u.Email == email {
			return converter.ToUserFromRepo(u), nil
		}
	}

	return user.User{}, user.ErrUserNotFound
}

func (r *Repository) Save(_ context.Context, u user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.users == nil {
		r.users = make(map[uuid.UUID]repoModel.User)
	}

	if _, ok := r.users[u.ID]; ok {
		return fmt.Errorf("task already exists: %w", user.ErrFailedToStoreUser)
	}

	r.users[u.ID] = converter.ToRepoFromUser(u)

	return nil
}
