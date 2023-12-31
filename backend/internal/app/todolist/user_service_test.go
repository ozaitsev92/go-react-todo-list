package todolist_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store/teststore"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/todolist"
	"github.com/ozaitsev92/go-react-todo-list/internal/helpers"
	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateTest(t *testing.T) {
	testCases := []struct {
		name    string
		r       func() *todolist.CreateUserRequest
		isValid bool
	}{
		{
			name: "valid",
			r: func() *todolist.CreateUserRequest {
				return &todolist.CreateUserRequest{
					Email:    "example@email.com",
					Password: "a password",
				}
			},
			isValid: true,
		},
		{
			name: "email is not present",
			r: func() *todolist.CreateUserRequest {
				return &todolist.CreateUserRequest{
					Email:    "",
					Password: "a password",
				}
			},
			isValid: false,
		},
		{
			name: "email is invalid",
			r: func() *todolist.CreateUserRequest {
				return &todolist.CreateUserRequest{
					Email:    "this is not an email",
					Password: "a password",
				}
			},
			isValid: false,
		},
		{
			name: "password is too short",
			r: func() *todolist.CreateUserRequest {
				return &todolist.CreateUserRequest{
					Email:    "example@email.com",
					Password: helpers.RandomString(3),
				}
			},
			isValid: false,
		},
		{
			name: "password is too long",
			r: func() *todolist.CreateUserRequest {
				return &todolist.CreateUserRequest{
					Email:    "example@email.com",
					Password: helpers.RandomString(101),
				}
			},
			isValid: false,
		},
	}

	service := todolist.NewUserService(teststore.New().User())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := service.CreateUser(context.Background(), tc.r())

			if tc.isValid {
				assert.NoError(t, err)
				assert.NotNil(t, u)
			} else {
				assert.Error(t, err)
				assert.Nil(t, u)
			}
		})
	}
}

func TestUserService_AuthenticateUser(t *testing.T) {
	testCases := []struct {
		name            string
		u               func() *todolist.User
		pwd             string
		isAuthenticated bool
	}{
		{
			name: "is authenticated",
			u: func() *todolist.User {
				return todolist.TestUser(t, "example@email.com", "a password")
			},
			pwd:             "a password",
			isAuthenticated: true,
		},
		{
			name: "incorrect password",
			u: func() *todolist.User {
				return todolist.TestUser(t, "example@email.com", "a password")
			},
			pwd:             "different password",
			isAuthenticated: false,
		},
	}

	store := teststore.New().User()
	service := todolist.NewUserService(store)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := tc.u()
			err := store.SaveUser(context.Background(), u)
			assert.NoError(t, err)

			r := &todolist.AuthenticateUserRequest{
				Email:    u.GetEmail(),
				Password: tc.pwd,
			}
			authenticatedUser, err := service.AuthenticateUser(context.Background(), r)

			if tc.isAuthenticated {
				assert.NoError(t, err)
				assert.NotNil(t, authenticatedUser)
			} else {
				assert.Error(t, err)
				assert.Nil(t, authenticatedUser)
			}
		})
	}
}

func TestUserService_FindUserByID(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *todolist.User
		isFound bool
	}{
		{
			name: "found",
			u: func() *todolist.User {
				return todolist.TestUser(t, "example@email.com", "a password")
			},
			isFound: true,
		},
		{
			name: "not found",
			u: func() *todolist.User {
				return todolist.TestUser(t, "example@email.com", "a password")
			},
			isFound: false,
		},
	}

	store := teststore.New().User()
	service := todolist.NewUserService(store)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := tc.u()
			err := store.SaveUser(context.Background(), u)
			assert.NoError(t, err)

			if tc.isFound {
				r := &todolist.FindUserByIDRequest{ID: u.GetID()}
				foundUser, err := service.FindUserByID(context.Background(), r)

				assert.NoError(t, err)
				assert.NotNil(t, foundUser)
			} else {
				r := &todolist.FindUserByIDRequest{ID: uuid.New()}
				foundUser, err := service.FindUserByID(context.Background(), r)

				assert.Error(t, err)
				assert.Nil(t, foundUser)
			}
		})
	}
}
