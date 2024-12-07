package task

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidText   = errors.New("text is invalid")
	ErrInvalidUserID = errors.New("user id is invalid")
)

// Task is a representation of a task entity.
type Task struct {
	ID        uuid.UUID
	Text      string
	Completed bool
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewTask creates and returns a new Task.
func NewTask(text string, userId uuid.UUID) (Task, error) {
	if text == "" {
		return Task{}, ErrInvalidText
	}

	if userId == uuid.Nil {
		return Task{}, ErrInvalidUserID
	}

	currentTime := time.Now()

	return Task{
		ID:        uuid.New(),
		Text:      text,
		Completed: false,
		UserID:    userId,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}, nil
}

// SetText sets the Text field.
func (t *Task) SetText(text string) error {
	if text == "" {
		return ErrInvalidText
	}

	t.Text = text
	t.UpdatedAt = time.Now()

	return nil
}

// MarkCompleted marks task a completed.
func (t *Task) MarkCompleted() {
	t.Completed = true
	t.UpdatedAt = time.Now()
}

// MarkNotCompleted marks task a NOT completed.
func (t *Task) MarkNotCompleted() {
	t.Completed = false
	t.UpdatedAt = time.Now()
}
