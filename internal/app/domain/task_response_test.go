package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestTaskToResponse(t *testing.T) {
	testCase := struct {
		text   string
		order  int
		isDone bool
		userID uuid.UUID
	}{
		"test text",
		1,
		false,
		uuid.New(),
	}

	task := domain.TestTask(t, testCase.text, testCase.order, testCase.isDone, testCase.userID)
	r := domain.TaskToResponse(*task)
	assert.NotEmpty(t, r.ID)
	assert.Equal(t, testCase.text, r.TaskText)
	assert.Equal(t, testCase.order, r.TaskOrder)
	assert.Equal(t, testCase.isDone, r.IsDone)
	assert.Equal(t, testCase.userID, r.UserID)
	assert.NotEmpty(t, r.CreatedAt)
	assert.NotEmpty(t, r.UpdatedAt)
}
