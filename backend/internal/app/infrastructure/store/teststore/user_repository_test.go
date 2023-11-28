package teststore_test

import (
	"context"
	"testing"

	"github.com/ozaitsev92/go-react-todo-list/internal/app/domain"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_SaveUser(t *testing.T) {
	s := teststore.New()
	u := domain.TestUser(t, "email@example.com", "a password")
	assert.NoError(t, s.User().SaveUser(context.Background(), u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.New()

	email := "user@example.org"

	u, err := s.User().FindByEmail(context.Background(), email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, u)

	u = domain.TestUser(t, email, "a password")
	assert.NoError(t, s.User().SaveUser(context.Background(), u))

	u, err = s.User().FindByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, email, u.GetEmail())
}

func TestUserRepository_FindByID(t *testing.T) {
	s := teststore.New()

	u := domain.TestUser(t, "email@example.com", "a password")
	assert.NoError(t, s.User().SaveUser(context.Background(), u))

	id := u.GetID()

	u, err := s.User().FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, id, u.GetID())
}
