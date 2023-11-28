package teststore

import (
	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/todolist"
)

type Store struct {
	userRepository todolist.UserRepository
	taskRepository todolist.TaskRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() todolist.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
			users: make(map[uuid.UUID]*todolist.User),
		}
	}

	return s.userRepository
}

func (s *Store) Task() todolist.TaskRepository {
	if s.taskRepository == nil {
		s.taskRepository = &TaskRepository{
			store: s,
			tasks: make(map[uuid.UUID]*todolist.Task),
		}
	}

	return s.taskRepository
}
