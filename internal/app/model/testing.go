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
		TaskText:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		TaskOrder: 0,
		IsDone:    false,
		UserID:    u.ID,
	}
	return task
}
