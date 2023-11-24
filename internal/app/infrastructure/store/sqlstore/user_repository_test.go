package sqlstore_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/domain"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

	u := domain.TestUser(t, "example@email.com", "a password")
	assert.NoError(t, s.User().SaveUser(context.Background(), u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

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

func TestUserRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

	id := uuid.New()

	u, err := s.User().FindByID(context.Background(), id)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, u)

	u = domain.TestUser(t, "example@email.com", "a password")
	assert.NoError(t, s.User().SaveUser(context.Background(), u))

	id = u.GetID()

	u, err = s.User().FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, id, u.GetID())
}
