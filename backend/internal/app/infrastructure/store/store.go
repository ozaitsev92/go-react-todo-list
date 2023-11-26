package store

import (
	"github.com/ozaitsev92/go-react-todo-list/internal/app/domain"
)

type Store interface {
	User() domain.UserRepository
	Task() domain.TaskRepository
}
