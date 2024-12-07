package model

import (
	"time"
)

type Task struct {
	ID          string    `bson:"_id"`
	Text        string    `bson:"text"`
	Completed   bool      `bson:"completed"`
	IsCompleted bool      `bson:"is_completed"`
	UserID      string    `bson:"user_id"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}
