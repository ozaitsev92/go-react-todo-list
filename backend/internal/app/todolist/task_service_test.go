package todolist_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store/teststore"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/todolist"
	"github.com/stretchr/testify/assert"
)

func TestTaskService_CreateTask(t *testing.T) {
	testCases := []struct {
		name    string
		r       func() *todolist.CreateTaskRequest
		isValid bool
	}{
		{
			name: "valid",
			r: func() *todolist.CreateTaskRequest {
				return &todolist.CreateTaskRequest{
					TaskText:  "task text",
					TaskOrder: 0,
				}
			},
			isValid: true,
		},
		{
			name: "task text invalid",
			r: func() *todolist.CreateTaskRequest {
				return &todolist.CreateTaskRequest{
					TaskText:  "",
					TaskOrder: 0,
				}
			},
			isValid: false,
		},
		{
			name: "task order invalid",
			r: func() *todolist.CreateTaskRequest {
				return &todolist.CreateTaskRequest{
					TaskText:  "task text",
					TaskOrder: -10,
				}
			},
			isValid: false,
		},
	}

	service := todolist.NewTaskService(teststore.New().Task())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := service.CreateTask(context.Background(), tc.r())

			if tc.isValid {
				assert.NoError(t, err)
				assert.NotNil(t, u)
			} else {
				assert.Error(t, err)
				assert.Nil(t, u)
			}
		})
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	testCases := []struct {
		name    string
		r       func(taskID uuid.UUID) *todolist.UpdateTaskRequest
		isValid bool
	}{
		{
			name: "valid",
			r: func(taskID uuid.UUID) *todolist.UpdateTaskRequest {
				return &todolist.UpdateTaskRequest{
					ID:        taskID,
					TaskText:  "task text",
					TaskOrder: 0,
				}
			},
			isValid: true,
		},
		{
			name: "task text invalid",
			r: func(taskID uuid.UUID) *todolist.UpdateTaskRequest {
				return &todolist.UpdateTaskRequest{
					ID:        taskID,
					TaskText:  "",
					TaskOrder: 0,
				}
			},
			isValid: false,
		},
		{
			name: "task order invalid",
			r: func(taskID uuid.UUID) *todolist.UpdateTaskRequest {
				return &todolist.UpdateTaskRequest{
					ID:        taskID,
					TaskText:  "task text",
					TaskOrder: -10,
				}
			},
			isValid: false,
		},
	}

	store := teststore.New().Task()
	service := todolist.NewTaskService(store)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task := todolist.TestTask(t, "texst text", 1, false, uuid.New())
			assert.NoError(t, store.SaveTask(context.Background(), task))

			u, err := service.UpdateTask(context.Background(), tc.r(task.GetID()))

			if tc.isValid {
				assert.NoError(t, err)
				assert.NotNil(t, u)
			} else {
				assert.Error(t, err)
				assert.Nil(t, u)
			}
		})
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	task := todolist.TestTask(t, "test text", 0, false, uuid.New())

	store := teststore.New().Task()
	assert.NoError(t, store.SaveTask(context.Background(), task))

	service := todolist.NewTaskService(store)
	r := &todolist.DeleteTaskRequest{ID: task.GetID()}
	assert.NoError(t, service.DeleteTask(context.Background(), r))
}

func TestTaskService_MarkTaskDone(t *testing.T) {
	task := todolist.TestTask(t, "test text", 0, false, uuid.New())

	store := teststore.New().Task()
	assert.NoError(t, store.SaveTask(context.Background(), task))

	service := todolist.NewTaskService(store)

	r := &todolist.MarkTaskDoneRequest{ID: task.GetID()}
	task, err := service.MarkTaskDone(context.Background(), r)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.True(t, task.GetIsDone())
}

func TestTaskService_MarkTaskNotDone(t *testing.T) {
	task := todolist.TestTask(t, "test text", 0, true, uuid.New())

	store := teststore.New().Task()
	assert.NoError(t, store.SaveTask(context.Background(), task))

	service := todolist.NewTaskService(store)

	r := &todolist.MarkTaskNotDoneRequest{ID: task.GetID()}
	task, err := service.MarkTaskNotDone(context.Background(), r)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.False(t, task.GetIsDone())
}

func TestTaskService_GetAllByUser(t *testing.T) {
	userID := uuid.New()
	task := todolist.TestTask(t, "test text", 0, true, userID)

	store := teststore.New().Task()
	assert.NoError(t, store.SaveTask(context.Background(), task))

	service := todolist.NewTaskService(store)
	r := &todolist.GetTasksByUserRequest{UserID: userID}
	tasks, err := service.GetAllByUser(context.Background(), r)
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, tasks[0].GetID(), task.GetID())
}
