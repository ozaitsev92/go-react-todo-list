package teststore

import (
	"context"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/domain"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store"
)

type UserRepository struct {
	store *Store
	users map[uuid.UUID]*domain.User
}

func (r *UserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	r.users[user.GetID()] = user
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, u := range r.users {
		if u.GetEmail() == email {
			return u, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}
