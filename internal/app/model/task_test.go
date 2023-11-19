package model_test

import (
	"testing"

	"github.com/ozaitsev92/go-react-todo-list/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestTask_Validate(t *testing.T) {
	u := model.TestUser(t)
	u.ID = 1

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
			name: "task Text is empty",
			task: func() *model.Task {
				task := model.TestTask(t, u)
				task.Text = ""
				return task
			},
			isValid: false,
		},
		{
			name: "task UserID is empty",
			task: func() *model.Task {
				task := model.TestTask(t, u)
				task.UserID = 0
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
