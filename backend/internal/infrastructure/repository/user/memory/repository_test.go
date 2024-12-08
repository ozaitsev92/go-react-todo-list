package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ozaitsev92/tododdd/config"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	repository "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/memory"
)

func TestRepositoryGetByID(t *testing.T) {
	cfg := config.Config{}

	u, err := user.NewUser("test1@example.com", "Password123")
	if err != nil {
		t.Errorf("GetByID() failed to create a new user: err = '%v'", err)
	}

	// Check if a user exists in the DB: should fail
	r := repository.NewRepository(cfg)
	_, err = r.GetByID(context.Background(), u.ID)
	if !errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("GetByID() got = '%v', want = '%v'", err, user.ErrUserNotFound)
	}

	// Add the user to the users collection
	err = r.Save(context.Background(), u)
	if err != nil {
		t.Errorf("GetByID() failed to save a new user: err = '%v'", err)
	}

	// Check if a user exists in the DB: should succeed
	foundUser, err := r.GetByID(context.Background(), u.ID)
	if errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("GetByID() err = '%v', want = '%v'", err, nil)
	}

	if foundUser.ID != u.ID {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundUser.ID, u.ID)
	}

	if foundUser.Email != u.Email {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundUser.Email, u.Email)
	}
}

func TestRepositoryGetByEmail(t *testing.T) {
	cfg := config.Config{}

	u, err := user.NewUser("test2@example.com", "Password123")
	if err != nil {
		t.Errorf("GetByEmail()failed to create a new user: err = '%v'", err)
	}

	// Check if a user exists in the DB: should fail
	r := repository.NewRepository(cfg)
	_, err = r.GetByEmail(context.Background(), u.Email)
	if !errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("GetByEmail() got = '%v', want = '%v'", err, user.ErrUserNotFound)
	}

	// Add the user to the users collection
	err = r.Save(context.Background(), u)
	if err != nil {
		t.Errorf("GetByID() failed to save a new user: err = '%v'", err)
	}

	// Check if a user exists in the DB: should succeed
	foundUser, err := r.GetByEmail(context.Background(), u.Email)
	if errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("GetByID() err = '%v', want = '%v'", err, nil)
	}

	if foundUser.ID != u.ID {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundUser.ID, u.ID)
	}

	if foundUser.Email != u.Email {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundUser.Email, u.Email)
	}
}

func TestRepositorySave(t *testing.T) {
	cfg := config.Config{}

	u, err := user.NewUser("test3@example.com", "Password123")
	if err != nil {
		t.Errorf("Save() failed to create a new user: err = '%v'", err)
	}

	// Check if a user exists in the DB: should fail
	r := repository.NewRepository(cfg)
	_, err = r.GetByEmail(context.Background(), u.Email)
	if !errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("GetByEmail() got = '%v', want = '%v'", err, user.ErrUserNotFound)
	}

	// Save the user into the DB
	err = r.Save(context.Background(), u)
	if err != nil {
		t.Errorf("Save() err = '%v'", err)
	}

	// Check if a user exists in the DB: should succeed
	foundUser, err := r.GetByEmail(context.Background(), u.Email)
	if errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("GetByEmail() got = '%v', want = '%v'", err, user.ErrUserNotFound)
	}

	if foundUser.ID != u.ID {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundUser.ID, u.ID)
	}

	if foundUser.Email != u.Email {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundUser.Email, u.Email)
	}
}
