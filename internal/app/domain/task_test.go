package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestTask_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		task    func() *domain.Task
		isValid bool
	}{
		{
			name: "valid",
			task: func() *domain.Task {
				return domain.TestTask(
					t,
					"test text",
					0,
					false,
					uuid.New(),
				)
			},
			isValid: true,
		},
		{
			name: "task TaskText is empty",
			task: func() *domain.Task {
				return domain.TestTask(
					t,
					"",
					0,
					false,
					uuid.New(),
				)
			},
			isValid: false,
		},
		{
			name: "task TaskOrder is invalid",
			task: func() *domain.Task {
				return domain.TestTask(
					t,
					"test text",
					-1,
					false,
					uuid.New(),
				)
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

func TestTask_MarkDone(t *testing.T) {
	task := domain.TestTask(t, "test text", 0, false, uuid.New())
	task.MarkDone()
	assert.True(t, task.GetIsDone())
}

func TestTask_MarkNotDone(t *testing.T) {
	task := domain.TestTask(t, "test text", 0, true, uuid.New())
	task.MarkNotDone()
	assert.False(t, task.GetIsDone())
}

func TestTask_BeforeUpdate(t *testing.T) {
	task := domain.TestTask(t, "test text", 0, true, uuid.New())
	assert.Equal(t, task.GetUpdatedAt(), task.GetCreatedAt())
	task.BeforeUpdate()
	assert.NotEqual(t, task.GetUpdatedAt(), task.GetCreatedAt())
}