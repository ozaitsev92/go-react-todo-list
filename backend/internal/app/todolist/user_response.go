package todolist

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func UserToResponse(u User) UserResponse {
	return UserResponse{
		ID:        u.id,
		Email:     u.email,
		CreatedAt: u.createdAt,
		UpdatedAt: u.updatedAt,
	}
}
