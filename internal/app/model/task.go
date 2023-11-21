package model

import (
	"time"

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
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Task) BeforeCreate() error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}

	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	return nil
}

func (t *Task) BeforeUpdate() error {
	t.UpdatedAt = time.Now()

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
