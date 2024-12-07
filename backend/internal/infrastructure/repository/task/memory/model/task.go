package model

import (
	"time"
)

type Task struct {
	ID        string
	Text      string
	Completed bool
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
