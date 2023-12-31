package todolist_test

import (
	"testing"

	"github.com/ozaitsev92/go-react-todo-list/internal/app/todolist"
	"github.com/ozaitsev92/go-react-todo-list/internal/helpers"
	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *todolist.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *todolist.User {
				return todolist.TestUser(t, "example@email.com", "a password")
			},
			isValid: true,
		},
		{
			name: "email is not present",
			u: func() *todolist.User {
				return todolist.TestUser(t, "", "a password")
			},
			isValid: false,
		},
		{
			name: "email is invalid",
			u: func() *todolist.User {
				return todolist.TestUser(t, "this is not an email", "a password")
			},
			isValid: false,
		},
		{
			name: "password is not present",
			u: func() *todolist.User {
				return todolist.TestUser(t, "example@email.com", "")
			},
			isValid: false,
		},
		{
			name: "password is too short",
			u: func() *todolist.User {
				return todolist.TestUser(t, "example@email.com", helpers.RandomString(3))
			},
			isValid: false,
		},
		{
			name: "password is too long",
			u: func() *todolist.User {
				return todolist.TestUser(t, "example@email.com", helpers.RandomString(101))
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}

func TestUser_BeforeUpdate(t *testing.T) {
	u := todolist.TestUser(t, "test text", "test password")
	assert.Equal(t, u.GetUpdatedAt(), u.GetCreatedAt())
	assert.NoError(t, u.BeforeUpdate())
	assert.NotEqual(t, u.GetUpdatedAt(), u.GetCreatedAt())
}

func TestUser_ComparePassword(t *testing.T) {
	testCases := []struct {
		name   string
		u      func() *todolist.User
		pwd    string
		isSame bool
	}{
		{
			name: "valid",
			u: func() *todolist.User {
				return todolist.TestUser(t, "example@email.com", "a password")
			},
			pwd:    "a password",
			isSame: true,
		},
		{
			name: "email is not present",
			u: func() *todolist.User {
				return todolist.TestUser(t, "example@email.com", "a password")
			},
			pwd:    "this is a different password",
			isSame: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isSame {
				assert.True(t, tc.u().ComparePassword(tc.pwd))
			} else {
				assert.False(t, tc.u().ComparePassword(tc.pwd))
			}
		})
	}
}
