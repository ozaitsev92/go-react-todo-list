package domain

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type Task struct {
	id        uuid.UUID
	taskText  string
	taskOrder int
	isDone    bool
	userID    uuid.UUID
	createdAt time.Time
	updatedAt time.Time
}

func CreateTask(taskText string, taskOrder int, userID uuid.UUID) (*Task, error) {
	currTime := time.Now()

	t := &Task{
		id:        uuid.New(),
		taskText:  taskText,
		taskOrder: taskOrder,
		isDone:    false,
		userID:    userID,
		createdAt: currTime,
		updatedAt: currTime,
	}

	if err := t.Validate(); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Task) GetID() uuid.UUID {
	return t.id
}

func (t *Task) GetTaskText() string {
	return t.taskText
}

func (t *Task) GetTaskOrder() int {
	return t.taskOrder
}

func (t *Task) GetIsDone() bool {
	return t.isDone
}

func (t *Task) GetUserID() uuid.UUID {
	return t.userID
}

func (t *Task) GetCreatedAt() time.Time {
	return t.createdAt
}

func (t *Task) GetUpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Task) MarkDone() {
	t.isDone = true
}

func (t *Task) MarkNotDone() {
	t.isDone = false
}

func (t *Task) BeforeUpdate() error {
	t.updatedAt = time.Now()

	return nil
}

func (t *Task) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.taskText, validation.Required, validation.Length(1, 255)),
		validation.Field(&t.taskOrder, validation.Min(0)),
		validation.Field(&t.userID, validation.Required, is.UUID),
	)
}
