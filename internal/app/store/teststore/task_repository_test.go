package teststore_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/model"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestTaskRepository_Create(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	u.ID = uuid.New()
	task := model.TestTask(t, u)
	assert.NoError(t, s.Task().Create(task))
	assert.NotNil(t, u)
}

func TestTaskRepository_GetAllByUser(t *testing.T) {
	s := teststore.New()
	taskRepo := s.Task()

	u := model.TestUser(t)
	u.ID = uuid.New()

	task1 := model.TestTask(t, u)
	task1.TaskText = "this is a completed task"
	task1.IsDone = true
	task1.TaskOrder = 0
	taskRepo.Create(task1)

	task2 := model.TestTask(t, u)
	task2.IsDone = true
	task2.TaskOrder = 1
	taskRepo.Create(task2)

	userTasks, err := taskRepo.GetAllByUser(u.ID)
	assert.NoError(t, err)
	assert.Len(t, userTasks, 2)

	var userTask1 *model.Task
	for _, t := range userTasks {
		if t.ID == task1.ID {
			userTask1 = t
			break
		}
	}

	assert.NotNil(t, userTask1)
	assert.Equal(t, task1.ID, userTasks[0].ID)
	assert.Equal(t, task1.TaskText, userTasks[0].TaskText)
	assert.Equal(t, task1.IsDone, userTasks[0].IsDone)
	assert.Equal(t, task1.UserID, userTasks[0].UserID)

	var userTask2 *model.Task
	for _, t := range userTasks {
		if t.ID == task2.ID {
			userTask2 = t
			break
		}
	}

	assert.NotNil(t, userTask2)
	assert.Equal(t, task2.ID, userTasks[1].ID)
	assert.Equal(t, task2.TaskText, userTasks[1].TaskText)
	assert.Equal(t, task2.IsDone, userTasks[1].IsDone)
	assert.Equal(t, task2.UserID, userTasks[1].UserID)
}

func TestTaskRepository_MarkAsDone(t *testing.T) {
	s := teststore.New()
	taskRepo := s.Task()

	u := model.TestUser(t)
	u.ID = uuid.New()
	task := model.TestTask(t, u)

	assert.NoError(t, taskRepo.Create(task))
	assert.NotNil(t, task)
	assert.False(t, task.IsDone)

	task, err := taskRepo.MarkAsDone(task.ID)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.True(t, task.IsDone)
}

func TestTaskRepository_MarkAsNotDone(t *testing.T) {
	s := teststore.New()
	taskRepo := s.Task()

	u := model.TestUser(t)
	u.ID = uuid.New()
	task := model.TestTask(t, u)
	task.IsDone = true

	assert.NoError(t, taskRepo.Create(task))
	assert.NotNil(t, task)
	assert.True(t, task.IsDone)

	task, err := taskRepo.MarkAsNotDone(task.ID)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.False(t, task.IsDone)
}

func TestTaskRepository_Delete(t *testing.T) {
	s := teststore.New()
	taskRepo := s.Task()

	u := model.TestUser(t)
	u.ID = uuid.New()
	task := model.TestTask(t, u)
	task.IsDone = true

	assert.NoError(t, taskRepo.Create(task))
	assert.NotNil(t, task)
	assert.True(t, task.IsDone)

	err := taskRepo.Delete(task.ID)
	assert.NoError(t, err)

	userTasks, err := taskRepo.GetAllByUser(u.ID)
	assert.NoError(t, err)
	assert.Len(t, userTasks, 0)
}
