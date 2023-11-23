package domain_test

import (
	"testing"

	"github.com/ozaitsev92/go-react-todo-list/internal/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserToResponse(t *testing.T) {
	testCase := struct {
		email    string
		password string
	}{
		"example@email.com",
		"a password",
	}

	user := domain.TestUser(t, testCase.email, testCase.password)
	r := domain.UserToResponse(*user)
	assert.NotEmpty(t, r.ID)
	assert.Equal(t, testCase.email, r.Email)
	assert.NotEmpty(t, r.CreatedAt)
	assert.NotEmpty(t, r.UpdatedAt)
}
