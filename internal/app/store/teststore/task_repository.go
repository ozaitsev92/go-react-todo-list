package teststore

import (
	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/model"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/store"
)

type TaskRepository struct {
	store *Store
	tasks map[uuid.UUID]*model.Task
}

func (r *TaskRepository) Create(t *model.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}

	if err := t.BeforeCreate(); err != nil {
		return err
	}

	r.tasks[t.ID] = t

	return nil
}

func (r *TaskRepository) GetAllByUser(UserID uuid.UUID) ([]*model.Task, error) {
	tasks := []*model.Task{}
	for _, t := range r.tasks {
		if t.UserID == UserID {
			tasks = append(tasks, t)
		}
	}

	return tasks, nil
}

func (r *TaskRepository) MarkAsDone(ID uuid.UUID) (*model.Task, error) {
	if t, ok := r.tasks[ID]; ok {
		t.IsDone = true
		r.tasks[ID] = t
		return t, nil
	} else {
		return nil, store.ErrRecordNotFound
	}
}

func (r *TaskRepository) MarkAsNotDone(ID uuid.UUID) (*model.Task, error) {
	if t, ok := r.tasks[ID]; ok {
		t.IsDone = false
		r.tasks[ID] = t
		return t, nil
	} else {
		return nil, store.ErrRecordNotFound
	}
}

func (r *TaskRepository) Delete(ID uuid.UUID) error {
	if _, ok := r.tasks[ID]; ok {
		delete(r.tasks, ID)
		return nil
	} else {
		return store.ErrRecordNotFound
	}
}
