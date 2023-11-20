package sqlstore_test

import (
	"testing"

	"github.com/ozaitsev92/go-react-todo-list/internal/app/model"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestTaskRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "tasks")

	s := sqlstore.New(db)

	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)

	task := model.TestTask(t, u)
	assert.NoError(t, s.Task().Create(task))
	assert.NotNil(t, task)
}

func TestTaskRepository_GetAllByUser(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "tasks")

	s := sqlstore.New(db)

	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)

	task1 := model.TestTask(t, u)
	task1.TaskOrder = 1
	task1.TaskText = "this is a first task"
	task1.IsDone = true
	assert.NoError(t, s.Task().Create(task1))
	assert.NotNil(t, task1)

	task2 := model.TestTask(t, u)
	task2.TaskOrder = 0
	task2.TaskText = "this is a first task"
	assert.NoError(t, s.Task().Create(task2))
	assert.NotNil(t, task2)

	userTasks, err := s.Task().GetAllByUser(u.ID)
	assert.NoError(t, err)
	assert.Len(t, userTasks, 2)

	assert.Equal(t, userTasks[0].ID, task2.ID)
	assert.Equal(t, userTasks[0].TaskOrder, task2.TaskOrder)
	assert.Equal(t, userTasks[0].TaskText, task2.TaskText)
	assert.Equal(t, userTasks[0].UserID, task2.UserID)
	assert.Equal(t, userTasks[0].IsDone, task2.IsDone)

	assert.Equal(t, userTasks[1].ID, task1.ID)
	assert.Equal(t, userTasks[1].TaskOrder, task1.TaskOrder)
	assert.Equal(t, userTasks[1].TaskText, task1.TaskText)
	assert.Equal(t, userTasks[1].UserID, task1.UserID)
	assert.Equal(t, userTasks[1].IsDone, task1.IsDone)
}

func TestTaskRepository_MarkAsDone(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "tasks")

	s := sqlstore.New(db)

	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)

	task := model.TestTask(t, u)
	assert.NoError(t, s.Task().Create(task))
	assert.NotNil(t, task)
	assert.False(t, task.IsDone)

	task, err := s.Task().MarkAsDone(task.ID)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.True(t, task.IsDone)
}

func TestTaskRepository_MarkAsNotDone(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "tasks")

	s := sqlstore.New(db)

	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)

	task := model.TestTask(t, u)
	task.IsDone = true
	assert.NoError(t, s.Task().Create(task))
	assert.NotNil(t, task)
	assert.True(t, task.IsDone)

	task, err := s.Task().MarkAsNotDone(task.ID)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.False(t, task.IsDone)
}

func TestTaskRepository_Delete(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "tasks")

	s := sqlstore.New(db)

	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)

	task := model.TestTask(t, u)
	task.IsDone = true
	assert.NoError(t, s.Task().Create(task))
	assert.NotNil(t, task)
	assert.True(t, task.IsDone)

	err := s.Task().Delete(task.ID)
	assert.NoError(t, err)

	userTasks, err := s.Task().GetAllByUser(u.ID)
	assert.NoError(t, err)
	assert.Len(t, userTasks, 0)
}
