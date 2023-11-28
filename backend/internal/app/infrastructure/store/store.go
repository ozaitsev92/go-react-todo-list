package store

import "github.com/ozaitsev92/go-react-todo-list/internal/app/todolist"

type Store interface {
	User() todolist.UserRepository
	Task() todolist.TaskRepository
}
