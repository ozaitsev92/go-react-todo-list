package app_test

import (
	"testing"

	"github.com/ozaitsev92/go-react-todo-list/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	testCases := []struct {
		name   string
		length int
		want   int
	}{
		{
			name:   "0 chars rand string",
			length: 0,
			want:   0,
		},
		{
			name:   "10 chars rand string",
			length: 10,
			want:   10,
		},
		{
			name:   "99 chars rand string",
			length: 99,
			want:   99,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Len(t, app.RandomString(tc.length), tc.length)
		})
	}
}
