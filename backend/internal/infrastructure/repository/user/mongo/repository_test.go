//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/ozaitsev92/tododdd/internal/domain/user"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/repository/contract"
	repository "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/mongo"
	"github.com/ozaitsev92/tododdd/internal/testutil/mongotest"
)

func TestMain(m *testing.M) { mongotest.Main(m) }

func TestUserRepository_Contract(t *testing.T) {
	contract.UserRepository(t, func(t *testing.T) user.Repository {
		r := repository.NewRepository(mongotest.DB(t))
		if err := r.EnsureIndexes(context.Background()); err != nil {
			t.Fatalf("EnsureIndexes: %v", err)
		}
		return r
	})
}
