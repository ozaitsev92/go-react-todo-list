package teststore

import (
	"context"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/todolist"
)

type UserRepository struct {
	store *Store
	users map[uuid.UUID]*todolist.User
}

func (r *UserRepository) SaveUser(ctx context.Context, user *todolist.User) error {
	r.users[user.GetID()] = user
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*todolist.User, error) {
	for _, u := range r.users {
		if u.GetEmail() == email {
			return u, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*todolist.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}
