package model

import (
	"time"
)

type Task struct {
	ID          string
	Text        string
	Completed   bool
	IsCompleted bool
	UserID      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
