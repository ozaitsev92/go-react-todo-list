package teststore

import (
	"github.com/ozaitsev92/go-react-todo-list/internal/app/model"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[string]*model.User
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	r.users[u.Email] = u
	u.ID = len(r.users)

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u, ok := r.users[email]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}

	return nil, store.ErrRecordNotFound
}
