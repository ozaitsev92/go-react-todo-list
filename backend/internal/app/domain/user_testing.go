package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUser(t *testing.T, email, password string) *User {
	var encryptedPassword string
	if len(password) > 0 {
		encryptedPassword, _ = encryptString(password)
	}

	currTime := time.Now()

	user := &User{
		id:                uuid.New(),
		email:             email,
		password:          password,
		encryptedPassword: encryptedPassword,
		createdAt:         currTime,
		updatedAt:         currTime,
	}

	return user
}
