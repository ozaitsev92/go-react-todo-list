package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/config"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	repository "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/memory"
)

func TestRepositoryGetByID(t *testing.T) {
	cfg := config.Config{}

	ti, err := task.NewTask("task text", uuid.New())
	if err != nil {
		t.Errorf("GetByID() failed to create a new task: err = '%v'", err)
	}

	// Check if a task exists in the DB: should fail
	r := repository.NewRepository(cfg)
	_, err = r.GetByID(context.Background(), ti.ID)
	if !errors.Is(err, task.ErrTaskNotFound) {
		t.Errorf("GetByID() got = '%v', want = '%v'", err, task.ErrTaskNotFound)
	}

	// Add the task to the tasks collection
	err = r.Save(context.Background(), ti)
	if err != nil {
		t.Errorf("GetByID() failed to save a new task: err = '%v'", err)
	}

	// Check if a task exists in the DB: should succeed
	foundTask, err := r.GetByID(context.Background(), ti.ID)
	if errors.Is(err, task.ErrTaskNotFound) {
		t.Errorf("GetByID() err = '%v', want = '%v'", err, nil)
	}

	if foundTask.ID != ti.ID {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundTask.ID, ti.ID)
	}

	if foundTask.Text != ti.Text {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundTask.Text, ti.Text)
	}

	if foundTask.UserID != ti.UserID {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundTask.UserID, ti.UserID)
	}
}

func TestRepositoryGetAllByUserID(t *testing.T) {
	cfg := config.Config{}

	userId1 := uuid.New()

	t1, err := task.NewTask("task text 1", userId1)
	if err != nil {
		t.Errorf("GetAllByUserID() failed to create a new task: err = '%v'", err)
	}

	t2, err := task.NewTask("task text 2", userId1)
	if err != nil {
		t.Errorf("GetAllByUserID() failed to create a new task: err = '%v'", err)
	}

	userId2 := uuid.New()

	t3, err := task.NewTask("task text 2", userId2)
	if err != nil {
		t.Errorf("GetAllByUserID() failed to create a new task: err = '%v'", err)
	}

	// Check if a task exists in the DB: should fail
	r := repository.NewRepository(cfg)
	foundTasks, err := r.GetAllByUserID(context.Background(), t1.UserID)
	if err != nil {
		t.Errorf("GetAllByUserID() error = '%v'", err)
	}

	if len(foundTasks) != 0 {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", len(foundTasks), 2)
	}

	// Add the task to the tasks collection
	for _, ti := range []task.Task{t1, t2, t3} {
		// Save the task into the DB
		err = r.Save(context.Background(), ti)
		if err != nil {
			t.Errorf("GetAllByUserID() failed to save new tasks: err = '%v'", err)
		}
	}

	// Check if a task exists in the DB: should succeed
	foundTasks, err = r.GetAllByUserID(context.Background(), userId1)
	if errors.Is(err, task.ErrTaskNotFound) {
		t.Errorf("GetAllByUserID() err = '%v', want = '%v'", err, nil)
	}

	if len(foundTasks) != 2 {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", len(foundTasks), 2)
	}

	// Task #1
	ft1 := foundTasks[0]
	if ft1.ID != t1.ID {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", ft1.ID, t1.ID)
	}

	if ft1.Text != t1.Text {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", ft1.Text, t1.Text)
	}

	if ft1.UserID != t1.UserID {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", ft1.UserID, t1.UserID)
	}

	// Task #2
	ft2 := foundTasks[1]
	if ft2.ID != t2.ID {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", ft2.ID, t2.ID)
	}

	if ft2.Text != t2.Text {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", ft2.Text, t2.Text)
	}

	if ft2.UserID != t2.UserID {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", ft2.UserID, t2.UserID)
	}
}

func TestRepositorySave(t *testing.T) {
	cfg := config.Config{}

	ti, err := task.NewTask("task text", uuid.New())
	if err != nil {
		t.Errorf("Save() failed to create a new task: err = '%v'", err)
	}

	// Check if a task exists in the DB: should fail
	r := repository.NewRepository(cfg)
	_, err = r.GetByID(context.Background(), ti.ID)
	if !errors.Is(err, task.ErrTaskNotFound) {
		t.Errorf("GetByID() got = '%v', want = '%v'", err, task.ErrTaskNotFound)
	}

	// Save the task into the DB
	err = r.Save(context.Background(), ti)
	if err != nil {
		t.Errorf("Save() err = '%v'", err)
	}

	// Check if a user exists in the DB: should succeed
	fountTask, err := r.GetByID(context.Background(), ti.ID)
	if errors.Is(err, task.ErrTaskNotFound) {
		t.Errorf("GetByID() got = '%v', want = '%v'", err, task.ErrTaskNotFound)
	}

	if fountTask.ID != ti.ID {
		t.Errorf("GetByID() got = '%v', want = '%v'", fountTask.ID, ti.ID)
	}

	if fountTask.UserID != ti.UserID {
		t.Errorf("GetByID() got = '%v', want = '%v'", fountTask.UserID, ti.UserID)
	}
}

func TestRepositoryUpdate(t *testing.T) {
	cfg := config.Config{}

	ti, err := task.NewTask("task text", uuid.New())
	if err != nil {
		t.Errorf("Update() failed to create a new task: err = '%v'", err)
	}

	// Add the task to the tasks collection
	r := repository.NewRepository(cfg)
	err = r.Save(context.Background(), ti)
	if err != nil {
		t.Errorf("Delete() failed to save a new task: err = '%v'", err)
	}

	ti.Text = "this is a new version of the task"

	// Update the task in the DB
	err = r.Update(context.Background(), ti)
	if err != nil {
		t.Errorf("Update() err = '%v'", err)
	}

	fountTask, err := r.GetByID(context.Background(), ti.ID)
	if errors.Is(err, task.ErrTaskNotFound) {
		t.Errorf("GetByID() err = '%v', want = '%v'", err, nil)
	}

	if fountTask.ID != ti.ID {
		t.Errorf("GetByID() got = '%v', want = '%v'", fountTask.ID, ti.ID)
	}

	if fountTask.UserID != ti.UserID {
		t.Errorf("GetByID() got = '%v', want = '%v'", fountTask.UserID, ti.UserID)
	}

	if fountTask.Text != ti.Text {
		t.Errorf("GetByID() got = '%v', want = '%v'", fountTask.Text, ti.Text)
	}
}

func TestRepositoryDelete(t *testing.T) {
	cfg := config.Config{}

	ti, err := task.NewTask("task text", uuid.New())
	if err != nil {
		t.Errorf("Delete() failed to create a new task: err = '%v'", err)
	}

	// Add the task to the tasks collection
	r := repository.NewRepository(cfg)
	err = r.Save(context.Background(), ti)
	if err != nil {
		t.Errorf("Delete() failed to save a new task: err = '%v'", err)
	}

	// Delete the task from the the DB
	err = r.Delete(context.Background(), ti.ID)
	if err != nil {
		t.Errorf("Delete() err = '%v'", err)
	}

	// Check if a task exists in the DB: should fail
	_, err = r.GetByID(context.Background(), ti.ID)
	if !errors.Is(err, task.ErrTaskNotFound) {
		t.Errorf("GetByID() err = '%v', want = '%v'", err, nil)
	}
}
