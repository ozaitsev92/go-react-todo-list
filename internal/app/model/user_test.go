package model_test

import (
	"testing"

	"github.com/ozaitsev92/go-react-todo-list/internal/app"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *model.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *model.User {
				return model.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "email is not present",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Email = ""
				return u
			},
			isValid: false,
		},
		{
			name: "email is invalid",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Email = "this is not an email"
				return u
			},
			isValid: false,
		},
		{
			name: "password is not present",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "password is too short",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = app.RandomString(3)
				return u
			},
			isValid: false,
		},
		{
			name: "password is too long",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = app.RandomString(101)
				return u
			},
			isValid: false,
		},
		{
			name: "valid if encrypted password is present",
			u: func() *model.User {
				u := model.TestUser(t)
				u.EncryptedPassword = "encrypted"
				u.Password = ""
				return u
			},
			isValid: true,
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

func TestUser_BeforeCreate(t *testing.T) {
	u := model.TestUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.ID)
	assert.NotEmpty(t, u.EncryptedPassword)
	assert.NotEqual(t, u.Password, u.EncryptedPassword)
}
