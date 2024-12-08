package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ozaitsev92/tododdd/config"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	repo "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/memory"
	"github.com/ozaitsev92/tododdd/internal/usecase"
)

func TestRegisterNewUser(t *testing.T) {
	type args struct {
		email    string
		password string
	}

	type testCase struct {
		name    string
		args    args
		want    user.User
		wantErr error
	}

	tests := []testCase{
		{
			name: "Success",
			args: args{
				email:    "test@example.com",
				password: "TestPassword1",
			},
			want: user.User{
				Email:    "test@example.com",
				Password: "TestPassword1",
			},
			wantErr: nil,
		},
		{
			name: "Empty email",
			args: args{
				email:    "",
				password: "TestPassword1",
			},
			want:    user.User{},
			wantErr: user.ErrInvalidEmail,
		},
		{
			name: "Empty password",
			args: args{
				email:    "test@example.com",
				password: "",
			},
			want:    user.User{},
			wantErr: user.ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := repo.NewRepository(config.Config{})

			s := usecase.NewUserUseCase(repo)
			newUser, err := s.RegisterNewUser(context.Background(), tt.args.email, tt.args.password)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("s.RegisterNewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && newUser.Email != tt.args.email {
				t.Errorf("s.RegisterNewUser() Email = %v, want %v", newUser.Email, tt.want.Email)
			}
			if err == nil && !newUser.ComparePassword(tt.want.Password) {
				t.Errorf("s.RegisterNewUser() Completed = %v, want %v", newUser.Password, tt.want.Password)
			}
			if err == nil && newUser.UpdatedAt.IsZero() {
				t.Error("s.RegisterNewUser() UpdatedAt is zero")
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	type args struct {
		email    string
		password string
	}

	type testCase struct {
		name    string
		args    args
		want    user.User
		wantErr error
	}

	tests := []testCase{
		{
			name: "Success",
			args: args{
				email:    "test@example.com",
				password: "TestPassword1",
			},
			want: user.User{
				Email:    "test@example.com",
				Password: "TestPassword1",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			newUser, err := user.NewUser(tt.args.email, tt.args.password)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("s.RegisterNewUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			repo := repo.NewRepository(config.Config{})
			err = repo.Save(context.Background(), newUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("s.RegisterNewUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			s := usecase.NewUserUseCase(repo)
			foundUser, err := s.GetUserByID(context.Background(), newUser.ID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("s.RegisterNewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && foundUser.Email != tt.args.email {
				t.Errorf("s.RegisterNewUser() Email = %v, want %v", newUser.Email, tt.want.Email)
			}
			if err == nil && !foundUser.ComparePassword(tt.want.Password) {
				t.Errorf("s.RegisterNewUser() Completed = %v, want %v", newUser.Password, tt.want.Password)
			}
			if err == nil && foundUser.UpdatedAt.IsZero() {
				t.Error("s.RegisterNewUser() UpdatedAt is zero")
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	type args struct {
		email    string
		password string
	}

	type testCase struct {
		name    string
		args    args
		want    user.User
		wantErr error
	}

	tests := []testCase{
		{
			name: "Success",
			args: args{
				email:    "test@example.com",
				password: "TestPassword1",
			},
			want: user.User{
				Email:    "test@example.com",
				Password: "TestPassword1",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			newUser, err := user.NewUser(tt.args.email, tt.args.password)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("s.RegisterNewUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			repo := repo.NewRepository(config.Config{})
			err = repo.Save(context.Background(), newUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("s.RegisterNewUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			s := usecase.NewUserUseCase(repo)
			foundUser, err := s.GetUserByEmail(context.Background(), newUser.Email)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("s.RegisterNewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && foundUser.Email != tt.args.email {
				t.Errorf("s.RegisterNewUser() Email = %v, want %v", newUser.Email, tt.want.Email)
			}
			if err == nil && !foundUser.ComparePassword(tt.want.Password) {
				t.Errorf("s.RegisterNewUser() Completed = %v, want %v", newUser.Password, tt.want.Password)
			}
			if err == nil && foundUser.UpdatedAt.IsZero() {
				t.Error("s.RegisterNewUser() UpdatedAt is zero")
			}
		})
	}
}
