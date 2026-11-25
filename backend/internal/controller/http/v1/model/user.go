package model

import (
	"time"

	"github.com/ozaitsev92/tododdd/internal/domain/user"
)

// User is the JSON representation of a user. It never includes the password hash.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponseFromUser converts a domain user into its JSON representation.
func ToResponseFromUser(u user.User) User {
	return User{
		ID:        u.ID.String(),
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
