package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ozaitsev92/tododdd/internal/domain/user"
	usermemory "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/memory"
	"github.com/ozaitsev92/tododdd/internal/usecase"
)

const validPassword = "Secret123"

func TestUserUseCase_RegisterNewUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("stores normalised email and hash, never the password", func(t *testing.T) {
		t.Parallel()
		repo := usermemory.NewRepository()
		uc, _ := newUserUseCase(t, repo)

		u, err := uc.RegisterNewUser(ctx, " Alice@Example.com ", validPassword)
		if err != nil {
			t.Fatal(err)
		}
		if u.Email != "alice@example.com" || u.PasswordHash != "hashed:"+validPassword {
			t.Fatalf("got %+v", u)
		}
		if stored, err := repo.GetByEmail(ctx, "alice@example.com"); err != nil || stored.ID != u.ID {
			t.Fatalf("not stored: %+v, %v", stored, err)
		}
	})

	t.Run("validation", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name, email, password string
			wantErr               error
		}{
			{"invalid email", "nope", validPassword, user.ErrInvalidEmail},
			{"short password", "a@example.com", "short", user.ErrWeakPassword},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				uc, _ := newUserUseCase(t, usermemory.NewRepository())
				if _, err := uc.RegisterNewUser(ctx, tt.email, tt.password); !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, want %v", err, tt.wantErr)
				}
			})
		}
	})

	t.Run("duplicate email in any case is rejected", func(t *testing.T) {
		t.Parallel()
		uc, _ := newUserUseCase(t, usermemory.NewRepository())
		if _, err := uc.RegisterNewUser(ctx, "bob@example.com", validPassword); err != nil {
			t.Fatal(err)
		}
		for _, email := range []string{"bob@example.com", "BOB@example.com", " bob@EXAMPLE.com"} {
			if _, err := uc.RegisterNewUser(ctx, email, validPassword); !errors.Is(err, user.ErrEmailTaken) {
				t.Fatalf("RegisterNewUser(%q) err = %v, want ErrEmailTaken", email, err)
			}
		}
	})

	t.Run("lookup failure is not reported as duplicate", func(t *testing.T) {
		t.Parallel()
		uc, _ := newUserUseCase(t, &faultyUsers{Repository: usermemory.NewRepository(), failGetByEmail: true})
		_, err := uc.RegisterNewUser(ctx, "c@example.com", validPassword)
		if !errors.Is(err, errDB) || errors.Is(err, user.ErrEmailTaken) {
			t.Fatalf("err = %v, want wrapped errDB", err)
		}
	})

	t.Run("save failure propagates", func(t *testing.T) {
		t.Parallel()
		uc, _ := newUserUseCase(t, &faultyUsers{Repository: usermemory.NewRepository(), failSave: true})
		if _, err := uc.RegisterNewUser(ctx, "d@example.com", validPassword); !errors.Is(err, errDB) {
			t.Fatalf("err = %v, want errDB", err)
		}
	})
}

func TestUserUseCase_Authenticate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	setup := func(t *testing.T) (*usecase.UserUseCase, *fakeHasher, user.User) {
		t.Helper()
		uc, h := newUserUseCase(t, usermemory.NewRepository())
		u, err := uc.RegisterNewUser(ctx, "alice@example.com", validPassword)
		if err != nil {
			t.Fatal(err)
		}
		h.compares = 0
		return uc, h, u
	}

	t.Run("correct credentials, email case-insensitive", func(t *testing.T) {
		t.Parallel()
		uc, _, want := setup(t)
		got, err := uc.Authenticate(ctx, "ALICE@example.com", validPassword)
		if err != nil || got.ID != want.ID {
			t.Fatalf("got %+v, %v", got, err)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		t.Parallel()
		uc, _, _ := setup(t)
		if _, err := uc.Authenticate(ctx, "alice@example.com", "wrong-password"); !errors.Is(err, usecase.ErrInvalidCredentials) {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("unknown and malformed email still run a hash comparison", func(t *testing.T) {
		t.Parallel()
		for _, email := range []string{"ghost@example.com", "not-an-email"} {
			uc, h, _ := setup(t)
			if _, err := uc.Authenticate(ctx, email, validPassword); !errors.Is(err, usecase.ErrInvalidCredentials) {
				t.Fatalf("Authenticate(%q) err = %v", email, err)
			}
			if h.compares != 1 {
				t.Fatalf("Authenticate(%q) ran %d comparisons, want 1 (timing-safe)", email, h.compares)
			}
		}
	})

	t.Run("lookup failure is an error, not invalid credentials", func(t *testing.T) {
		t.Parallel()
		uc, _ := newUserUseCase(t, &faultyUsers{Repository: usermemory.NewRepository(), failGetByEmail: true})
		_, err := uc.Authenticate(ctx, "alice@example.com", validPassword)
		if !errors.Is(err, errDB) {
			t.Fatalf("err = %v, want errDB", err)
		}
	})
}
