// Package repository is a concurrency-safe in-memory user.Repository.
package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
)

var _ user.Repository = (*Repository)(nil)

// Repository stores users by ID with a unique email index.
type Repository struct {
	mu      sync.RWMutex
	users   map[uuid.UUID]user.User
	byEmail map[string]uuid.UUID
}

// NewRepository returns an empty repository.
func NewRepository() *Repository {
	return &Repository{
		users:   make(map[uuid.UUID]user.User),
		byEmail: make(map[string]uuid.UUID),
	}
}

// GetByID returns the user with id or user.ErrUserNotFound.
func (r *Repository) GetByID(_ context.Context, id uuid.UUID) (user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return user.User{}, user.ErrUserNotFound
	}

	return u, nil
}

// GetByEmail returns the user with the (normalised) email or user.ErrUserNotFound.
func (r *Repository) GetByEmail(_ context.Context, email string) (user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byEmail[email]
	if !ok {
		return user.User{}, user.ErrUserNotFound
	}

	return r.users[id], nil
}

// Save inserts a user, enforcing unique IDs and emails.
func (r *Repository) Save(_ context.Context, u user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byEmail[u.Email]; ok {
		return fmt.Errorf("save user: %w", user.ErrEmailTaken)
	}
	if _, ok := r.users[u.ID]; ok {
		return fmt.Errorf("save user: duplicate id %s", u.ID)
	}

	r.users[u.ID] = u
	r.byEmail[u.Email] = u.ID

	return nil
}
