package store

import (
	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/model"
)

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	Find(uuid.UUID) (*model.User, error)
}

type TaskRepository interface {
	Create(*model.Task) error
	GetAllByUser(uuid.UUID) ([]*model.Task, error)
	MarkAsDone(uuid.UUID) (*model.Task, error)
	MarkAsNotDone(uuid.UUID) (*model.Task, error)
	Delete(uuid.UUID) error
}
