package model

import "testing"

func TestUser(t *testing.T) *User {
	user := &User{
		Email:    "user@example.org",
		Password: "password",
	}
	return user
}

func TestTask(t *testing.T, u *User) *Task {
	task := &Task{
		Text:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		UserID: u.ID,
	}
	return task
}
