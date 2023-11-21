package sqlstore_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/model"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/store"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

	email := "user@example.org"

	u, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, u)

	u = model.TestUser(t)
	u.Email = email
	s.User().Create(u)

	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, email, u.Email)
}

func TestUserRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

	id := uuid.New()

	u, err := s.User().Find(id)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, u)

	u = model.TestUser(t)
	s.User().Create(u)

	id = u.ID

	u, err = s.User().Find(id)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, id, u.ID)
}
