package model_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestTask_Validate(t *testing.T) {
	u := model.TestUser(t)
	u.ID = uuid.New()

	testCases := []struct {
		name    string
		task    func() *model.Task
		isValid bool
	}{
		{
			name: "valid",
			task: func() *model.Task {
				return model.TestTask(t, u)
			},
			isValid: true,
		},
		{
			name: "task TaskText is empty",
			task: func() *model.Task {
				task := model.TestTask(t, u)
				task.TaskText = ""
				return task
			},
			isValid: false,
		},
		{
			name: "task TaskOrder is invalid",
			task: func() *model.Task {
				task := model.TestTask(t, u)
				task.TaskOrder = -1
				return task
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.task().Validate())
			} else {
				assert.Error(t, tc.task().Validate())
			}
		})
	}
}

func TestTask_BeforeCreate(t *testing.T) {
	u := model.TestUser(t)
	u.ID = uuid.New()

	task := model.TestTask(t, u)
	assert.NoError(t, task.BeforeCreate())
	assert.NotEmpty(t, task.ID)
}
