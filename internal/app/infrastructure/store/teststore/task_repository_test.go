package teststore_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/domain"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestTaskRepository_SaveTask(t *testing.T) {
	s := teststore.New()
	task := domain.TestTask(t, "test task", 1, false, uuid.New())
	assert.NoError(t, s.Task().SaveTask(context.Background(), task))
	assert.NotNil(t, task)
}

func TestTaskRepository_FindByID(t *testing.T) {
	s := teststore.New()
	taskRepo := s.Task()

	task := domain.TestTask(t, "test task", 1, false, uuid.New())
	assert.NoError(t, taskRepo.SaveTask(context.Background(), task))
	assert.NotNil(t, task)

	taskFound, err := taskRepo.FindByID(context.Background(), task.GetID())
	assert.NoError(t, err)
	assert.NotNil(t, taskFound)
	assert.Equal(t, taskFound.GetID(), task.GetID())
}

func TestTaskRepository_GetAllByUserID(t *testing.T) {
	s := teststore.New()
	taskRepo := s.Task()

	userID := uuid.New()

	task1 := domain.TestTask(t, "this is a completed task", 0, true, userID)
	assert.NoError(t, taskRepo.SaveTask(context.Background(), task1))

	task2 := domain.TestTask(t, "this is a pending task", 1, false, userID)
	assert.NoError(t, taskRepo.SaveTask(context.Background(), task2))

	userTasks, err := taskRepo.GetAllByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, userTasks, 2)

	var userTask1 *domain.Task
	for _, t := range userTasks {
		if t.GetID() == task1.GetID() {
			userTask1 = t
			break
		}
	}

	assert.NotNil(t, userTask1)
	assert.Equal(t, task1.GetID(), userTasks[0].GetID())
	assert.Equal(t, task1.GetTaskText(), userTasks[0].GetTaskText())
	assert.Equal(t, task1.GetIsDone(), userTasks[0].GetIsDone())
	assert.Equal(t, task1.GetUserID(), userTasks[0].GetUserID())

	var userTask2 *domain.Task
	for _, t := range userTasks {
		if t.GetID() == task2.GetID() {
			userTask2 = t
			break
		}
	}

	assert.NotNil(t, userTask2)
	assert.Equal(t, task2.GetID(), userTasks[1].GetID())
	assert.Equal(t, task2.GetTaskText(), userTasks[1].GetTaskText())
	assert.Equal(t, task2.GetIsDone(), userTasks[1].GetIsDone())
	assert.Equal(t, task2.GetUserID(), userTasks[1].GetUserID())
}

func TestTaskRepository_DeleteTask(t *testing.T) {
	s := teststore.New()
	taskRepo := s.Task()

	userID := uuid.New()
	task := domain.TestTask(t, "test task", 1, false, userID)

	assert.NoError(t, taskRepo.SaveTask(context.Background(), task))
	assert.NotNil(t, task)

	err := taskRepo.DeleteTask(context.Background(), task)
	assert.NoError(t, err)

	userTasks, err := taskRepo.GetAllByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, userTasks, 0)
}
