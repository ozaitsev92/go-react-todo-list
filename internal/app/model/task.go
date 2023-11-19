package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Task struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Order  int    `json:"order"`
	IsDone bool   `json:"is_done"`
	UserID int    `json:"user_id"`
}

func (t *Task) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.Text, validation.Required, validation.Length(1, 255)),
		validation.Field(&t.UserID, validation.Required),
	)
}
