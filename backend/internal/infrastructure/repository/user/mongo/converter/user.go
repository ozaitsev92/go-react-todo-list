// Package converter maps between user.User and its MongoDB document.
package converter

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	repoModel "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/mongo/model"
)

// ToUserFromRepo converts a stored document into a domain user.
func ToUserFromRepo(u repoModel.User) (user.User, error) {
	id, err := uuid.Parse(u.ID)
	if err != nil {
		return user.User{}, fmt.Errorf("users: invalid _id %q: %w", u.ID, err)
	}

	return user.User{
		ID:           id,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt.UTC(),
		UpdatedAt:    u.UpdatedAt.UTC(),
	}, nil
}

// ToRepoFromUser converts a domain user into its stored document.
func ToRepoFromUser(u user.User) repoModel.User {
	return repoModel.User{
		ID:           u.ID.String(),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
