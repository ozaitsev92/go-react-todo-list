package todolist

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTask(t *testing.T, text string, order int, isDone bool, userID uuid.UUID) *Task {
	currTime := time.Now()

	task := &Task{
		id:        uuid.New(),
		taskText:  text,
		taskOrder: order,
		isDone:    isDone,
		userID:    userID,
		createdAt: currTime,
		updatedAt: currTime,
	}

	return task
}
