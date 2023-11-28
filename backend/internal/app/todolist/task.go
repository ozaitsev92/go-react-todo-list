package todolist

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

func (t *Task) SetID(id uuid.UUID) error {
	t.id = id
	return nil
}

func (t *Task) GetTaskText() string {
	return t.taskText
}

func (t *Task) SetTaskText(taskText string) error {
	err := validation.Validate(taskText, validation.Required, validation.Length(1, 255))
	if err != nil {
		return err
	}
	t.taskText = taskText
	return nil
}

func (t *Task) GetTaskOrder() int {
	return t.taskOrder
}

func (t *Task) SetTaskOrder(taskOrder int) error {
	err := validation.Validate(taskOrder, validation.Min(0))
	if err != nil {
		return err
	}
	t.taskOrder = taskOrder
	return nil
}

func (t *Task) GetIsDone() bool {
	return t.isDone
}

func (t *Task) SetIsDone(isDone bool) error {
	t.isDone = isDone
	return nil
}

func (t *Task) GetUserID() uuid.UUID {
	return t.userID
}

func (t *Task) SetUserID(userID uuid.UUID) error {
	t.userID = userID
	return nil
}

func (t *Task) GetCreatedAt() time.Time {
	return t.createdAt
}

func (t *Task) SetCreatedAt(createdAt time.Time) error {
	err := validation.Validate(createdAt, validation.By(timeNotZero))
	if err != nil {
		return err
	}

	t.createdAt = createdAt
	return nil
}

func (t *Task) GetUpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Task) SetUpdatedAt(updatedAt time.Time) error {
	err := validation.Validate(updatedAt, validation.By(timeNotZero))
	if err != nil {
		return err
	}

	t.updatedAt = updatedAt
	return nil
}

func (t *Task) MarkDone() error {
	return t.SetIsDone(true)
}

func (t *Task) MarkNotDone() error {
	return t.SetIsDone(false)
}

func (t *Task) BeforeUpdate() error {
	return t.SetUpdatedAt(time.Now())
}

func (t *Task) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.taskText, validation.Required, validation.Length(1, 255)),
		validation.Field(&t.taskOrder, validation.Min(0)),
		validation.Field(&t.userID, validation.Required, is.UUID),
		validation.Field(&t.createdAt, validation.By(timeNotZero)),
		validation.Field(&t.updatedAt, validation.By(timeNotZero)),
	)
}
