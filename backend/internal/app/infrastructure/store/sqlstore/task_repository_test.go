package sqlstore_test

import (
	"context"
	"testing"

	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store/sqlstore"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/todolist"
	"github.com/stretchr/testify/assert"
)

func TestTaskRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "tasks")

	s := sqlstore.New(db)

	u := todolist.TestUser(t, "example@email.com", "a password")
	assert.NoError(t, s.User().SaveUser(context.Background(), u))
	assert.NotNil(t, u)

	task := todolist.TestTask(t, "this is a first task", 0, false, u.GetID())
	assert.NoError(t, s.Task().SaveTask(context.Background(), task))
	assert.NotNil(t, task)
}

func TestTaskRepository_FindByID(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "tasks")

	s := sqlstore.New(db)

	u := todolist.TestUser(t, "example@email.com", "a password")
	assert.NoError(t, s.User().SaveUser(context.Background(), u))
	assert.NotNil(t, u)

	task := todolist.TestTask(t, "this is a first task", 0, true, u.GetID())
	assert.NoError(t, s.Task().SaveTask(context.Background(), task))
	assert.NotNil(t, task)

	foundTask, err := s.Task().FindByID(context.Background(), task.GetID())
	assert.NoError(t, err)
	assert.NotNil(t, foundTask)
	assert.Equal(t, task.GetID(), foundTask.GetID())
}

func TestTaskRepository_GetAllByUserID(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "tasks")

	s := sqlstore.New(db)

	u := todolist.TestUser(t, "example@email.com", "a password")
	assert.NoError(t, s.User().SaveUser(context.Background(), u))
	assert.NotNil(t, u)

	task1 := todolist.TestTask(t, "this is a first task", 0, true, u.GetID())
	assert.NoError(t, s.Task().SaveTask(context.Background(), task1))
	assert.NotNil(t, task1)

	task2 := todolist.TestTask(t, "this is a second task", 1, true, u.GetID())
	assert.NoError(t, s.Task().SaveTask(context.Background(), task2))
	assert.NotNil(t, task2)

	userTasks, err := s.Task().GetAllByUserID(context.Background(), u.GetID())
	assert.NoError(t, err)
	assert.Len(t, userTasks, 2)

	assert.Equal(t, userTasks[0].GetID(), task1.GetID())
	assert.Equal(t, userTasks[0].GetTaskOrder(), task1.GetTaskOrder())
	assert.Equal(t, userTasks[0].GetTaskText(), task1.GetTaskText())
	assert.Equal(t, userTasks[0].GetUserID(), task1.GetUserID())
	assert.Equal(t, userTasks[0].GetIsDone(), task1.GetIsDone())

	assert.Equal(t, userTasks[1].GetID(), task2.GetID())
	assert.Equal(t, userTasks[1].GetTaskOrder(), task2.GetTaskOrder())
	assert.Equal(t, userTasks[1].GetTaskText(), task2.GetTaskText())
	assert.Equal(t, userTasks[1].GetUserID(), task2.GetUserID())
	assert.Equal(t, userTasks[1].GetIsDone(), task2.GetIsDone())
}

func TestTaskRepository_Delete(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "tasks")

	s := sqlstore.New(db)

	u := todolist.TestUser(t, "example@email.com", "a password")
	assert.NoError(t, s.User().SaveUser(context.Background(), u))
	assert.NotNil(t, u)

	task := todolist.TestTask(t, "text", 0, true, u.GetID())
	assert.NoError(t, s.Task().SaveTask(context.Background(), task))
	assert.NotNil(t, task)
	assert.True(t, task.GetIsDone())

	err := s.Task().DeleteTask(context.Background(), task)
	assert.NoError(t, err)

	userTasks, err := s.Task().GetAllByUserID(context.Background(), u.GetID())
	assert.NoError(t, err)
	assert.Len(t, userTasks, 0)
}
