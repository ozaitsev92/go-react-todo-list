package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// TODO replace ID int with UUID
type Task struct {
	ID        int    `json:"id"`
	TaskText  string `json:"task_text"`
	TaskOrder int    `json:"order"`
	IsDone    bool   `json:"is_done"`
	UserID    int    `json:"user_id"`
}

func (t *Task) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.TaskText, validation.Required, validation.Length(1, 255)),
		validation.Field(&t.UserID, validation.Required),
	)
}
