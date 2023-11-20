package store

import "github.com/ozaitsev92/go-react-todo-list/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	Find(int) (*model.User, error)
}

type TaskRepository interface {
	Create(*model.Task) error
	GetAllByUser(int) ([]*model.Task, error)
	MarkAsDone(int) (*model.Task, error)
	MarkAsNotDone(int) (*model.Task, error)
	Delete(int) error
}
