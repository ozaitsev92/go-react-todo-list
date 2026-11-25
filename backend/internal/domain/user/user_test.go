package user_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
)

func TestNormalizeEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{in: "alice@example.com", want: "alice@example.com"},
		{in: "  Alice@Example.COM ", want: "alice@example.com"},
		{in: "", wantErr: true},
		{in: "not-an-email", wantErr: true},
		{in: "Alice <alice@example.com>", wantErr: true},
		{in: "<alice@example.com>", wantErr: true},
		{in: "a@b.c, d@e.f", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			got, err := user.NormalizeEmail(tt.in)
			if tt.wantErr {
				if !errors.Is(err, user.ErrInvalidEmail) {
					t.Fatalf("err = %v, want ErrInvalidEmail", err)
				}
				return
			}
			if err != nil || got != tt.want {
				t.Fatalf("NormalizeEmail(%q) = %q, %v; want %q", tt.in, got, err, tt.want)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		pw   string
		ok   bool
	}{
		{"min length", strings.Repeat("a", user.MinPasswordLength), true},
		{"max length", strings.Repeat("a", user.MaxPasswordLength), true},
		{"too short", strings.Repeat("a", user.MinPasswordLength-1), false},
		{"too long for bcrypt", strings.Repeat("a", user.MaxPasswordLength+1), false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := user.ValidatePassword(tt.pw)
			if tt.ok != (err == nil) || (err != nil && !errors.Is(err, user.ErrWeakPassword)) {
				t.Fatalf("ValidatePassword(len %d) = %v", len(tt.pw), err)
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	t.Parallel()

	t.Run("valid normalises email", func(t *testing.T) {
		t.Parallel()
		u, err := user.NewUser(" Bob@Example.com", "hash")
		if err != nil {
			t.Fatal(err)
		}
		if u.ID == uuid.Nil || u.Email != "bob@example.com" || u.PasswordHash != "hash" || u.CreatedAt.IsZero() {
			t.Fatalf("unexpected user %+v", u)
		}
	})

	t.Run("missing hash", func(t *testing.T) {
		t.Parallel()
		if _, err := user.NewUser("a@example.com", ""); !errors.Is(err, user.ErrInvalidPasswordHash) {
			t.Fatalf("err = %v", err)
		}
	})
}
