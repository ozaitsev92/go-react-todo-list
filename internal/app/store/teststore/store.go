package teststore

import (
	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/model"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/store"
)

type Store struct {
	userRepository *UserRepository
	taskRepository *TaskRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
			users: make(map[string]*model.User),
		}
	}

	return s.userRepository
}

func (s *Store) Task() store.TaskRepository {
	if s.userRepository == nil {
		s.taskRepository = &TaskRepository{
			store: s,
			tasks: make(map[uuid.UUID]*model.Task),
		}
	}

	return s.taskRepository
}
