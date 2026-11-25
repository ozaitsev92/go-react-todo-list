package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	taskmemory "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/memory"
	usermemory "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/memory"
	"github.com/ozaitsev92/tododdd/internal/usecase"
)

var errDB = errors.New("db unavailable")

// fakeHasher is a fast, deterministic PasswordHasher that counts Compare calls.
type fakeHasher struct{ compares int }

func (h *fakeHasher) Hash(pw string) (string, error) { return "hashed:" + pw, nil }

func (h *fakeHasher) Compare(hash, pw string) error {
	h.compares++
	if hash != "hashed:"+pw {
		return errors.New("mismatch")
	}
	return nil
}

// faultyTasks wraps the memory repository and fails selected operations.
type faultyTasks struct {
	*taskmemory.Repository
	failGet, failUpdate, failDelete, failSave bool
}

func (f *faultyTasks) GetByID(ctx context.Context, id uuid.UUID) (task.Task, error) {
	if f.failGet {
		return task.Task{}, errDB
	}
	return f.Repository.GetByID(ctx, id)
}

func (f *faultyTasks) Save(ctx context.Context, t task.Task) error {
	if f.failSave {
		return errDB
	}
	return f.Repository.Save(ctx, t)
}

func (f *faultyTasks) Update(ctx context.Context, t task.Task) error {
	if f.failUpdate {
		return errDB
	}
	return f.Repository.Update(ctx, t)
}

func (f *faultyTasks) Delete(ctx context.Context, id uuid.UUID) error {
	if f.failDelete {
		return errDB
	}
	return f.Repository.Delete(ctx, id)
}

// faultyUsers wraps the memory repository and fails selected operations.
type faultyUsers struct {
	*usermemory.Repository
	failGetByEmail, failSave bool
}

func (f *faultyUsers) GetByEmail(ctx context.Context, email string) (user.User, error) {
	if f.failGetByEmail {
		return user.User{}, errDB
	}
	return f.Repository.GetByEmail(ctx, email)
}

func (f *faultyUsers) Save(ctx context.Context, u user.User) error {
	if f.failSave {
		return errDB
	}
	return f.Repository.Save(ctx, u)
}

func newUserUseCase(t *testing.T, repo user.Repository) (*usecase.UserUseCase, *fakeHasher) {
	t.Helper()
	h := &fakeHasher{}
	uc, err := usecase.NewUserUseCase(repo, h)
	if err != nil {
		t.Fatal(err)
	}
	return uc, h
}
