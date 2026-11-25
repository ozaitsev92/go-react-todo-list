// Package task contains the task entity and its repository port.
package task

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// MaxTextLength is the maximum task text length in characters.
const MaxTextLength = 1000

var (
	// ErrInvalidText is returned for empty, whitespace-only or too long text.
	ErrInvalidText = errors.New("text must be 1 to 1000 characters")
	// ErrInvalidUserID is returned when a task has no owner.
	ErrInvalidUserID = errors.New("user id is invalid")
)

// Task is a to-do item owned by a single user.
type Task struct {
	ID        uuid.UUID
	Text      string
	Completed bool
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewTask creates a task for userID.
func NewTask(text string, userID uuid.UUID) (Task, error) {
	text, err := normalizeText(text)
	if err != nil {
		return Task{}, err
	}

	if userID == uuid.Nil {
		return Task{}, ErrInvalidUserID
	}

	now := now()

	return Task{
		ID:        uuid.New(),
		Text:      text,
		Completed: false,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// IsOwnedBy reports whether userID owns the task.
func (t *Task) IsOwnedBy(userID uuid.UUID) bool {
	return t.UserID == userID
}

// SetText replaces the text. The task is unchanged on error.
func (t *Task) SetText(text string) error {
	text, err := normalizeText(text)
	if err != nil {
		return err
	}

	t.Text = text
	t.UpdatedAt = now()

	return nil
}

// SetCompleted marks the task as completed or not completed.
func (t *Task) SetCompleted(completed bool) {
	t.Completed = completed
	t.UpdatedAt = now()
}

func normalizeText(text string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" || utf8.RuneCountInString(text) > MaxTextLength {
		return "", ErrInvalidText
	}

	return text, nil
}

// now returns UTC time truncated to MongoDB's millisecond precision.
func now() time.Time {
	return time.Now().UTC().Truncate(time.Millisecond)
}
