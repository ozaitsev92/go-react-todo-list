package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID `json:"id"`
	TaskText  string    `json:"task_text"`
	TaskOrder int       `json:"order"`
	IsDone    bool      `json:"is_done"`
	UserID    uuid.UUID `json:"user_id"`
}

func (t *Task) BeforeCreate() error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}

	return nil
}

func (t *Task) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.TaskText, validation.Required, validation.Length(1, 255)),
		validation.Field(&t.TaskOrder, validation.Min(0)),
		validation.Field(&t.UserID, validation.Required, is.UUID),
	)
}
