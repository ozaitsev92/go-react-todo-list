package converter

import (
	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	repoModel "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/memory/model"
)

func ToUserFromRepo(u repoModel.User) user.User {
	return user.User{
		ID:        uuid.MustParse(u.ID),
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToRepoFromUser(u user.User) repoModel.User {
	return repoModel.User{
		ID:        u.ID.String(),
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
