// Package model holds the MongoDB document shape for users.
package model

import (
	"time"
)

// User is the stored user document; "password" holds the bcrypt hash.
type User struct {
	ID           string    `bson:"_id"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}
