package repository_test

import (
	"testing"

	"github.com/ozaitsev92/tododdd/internal/domain/user"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/repository/contract"
	repository "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/memory"
)

func TestUserRepository_Contract(t *testing.T) {
	contract.UserRepository(t, func(*testing.T) user.Repository { return repository.NewRepository() })
}
