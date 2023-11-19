package model

import "testing"

func TestUser(t *testing.T) *User {
	u := &User{
		Email:    "user@example.org",
		Password: "password",
	}
	return u
}
