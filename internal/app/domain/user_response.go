package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UserToResponse(u User) UserResponse {
	return UserResponse{
		ID:        u.id,
		Email:     u.email,
		CreatedAt: u.createdAt,
		UpdatedAt: u.updatedAt,
	}
}
