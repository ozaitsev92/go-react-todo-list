package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/todolist"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
	taskRepository *TaskRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() todolist.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
		}
	}

	return s.userRepository
}

func (s *Store) Task() todolist.TaskRepository {
	if s.taskRepository == nil {
		s.taskRepository = &TaskRepository{
			store: s,
		}
	}

	return s.taskRepository
}
