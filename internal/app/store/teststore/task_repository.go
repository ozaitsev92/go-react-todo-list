package teststore

import (
	"github.com/ozaitsev92/go-react-todo-list/internal/app/model"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/store"
)

type TaskRepository struct {
	store *Store
	tasks map[int]*model.Task
}

func (r *TaskRepository) Create(t *model.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}

	t.ID = len(r.tasks) + 1
	r.tasks[t.ID] = t

	return nil
}

func (r *TaskRepository) GetAllByUser(UserID int) ([]*model.Task, error) {
	tasks := []*model.Task{}
	for _, t := range r.tasks {
		if t.UserID == UserID {
			tasks = append(tasks, t)
		}
	}

	return tasks, nil
}

func (r *TaskRepository) MarkAsDone(ID int) (*model.Task, error) {
	if t, ok := r.tasks[ID]; ok {
		t.IsDone = true
		r.tasks[ID] = t
		return t, nil
	} else {
		return nil, store.ErrRecordNotFound
	}
}

func (r *TaskRepository) MarkAsNotDone(ID int) (*model.Task, error) {
	if t, ok := r.tasks[ID]; ok {
		t.IsDone = false
		r.tasks[ID] = t
		return t, nil
	} else {
		return nil, store.ErrRecordNotFound
	}
}
