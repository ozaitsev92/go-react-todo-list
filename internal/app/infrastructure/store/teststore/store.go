package teststore

import (
	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/domain"
)

type Store struct {
	userRepository domain.UserRepository
	taskRepository domain.TaskRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() domain.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
			users: make(map[uuid.UUID]*domain.User),
		}
	}

	return s.userRepository
}

func (s *Store) Task() domain.TaskRepository {
	if s.taskRepository == nil {
		s.taskRepository = &TaskRepository{
			store: s,
			tasks: make(map[uuid.UUID]*domain.Task),
		}
	}

	return s.taskRepository
}
