package user_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

func TestUserNewUser(t *testing.T) {
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
				password: "test",
			},
			want: user.User{
				Email:    "test@example.com",
				Password: encryptString("test"),
			},
		},
		{
			name: "Empty email",
			args: args{
				email:    "",
				password: "test",
			},
			want:    user.User{},
			wantErr: user.ErrInvalidEmail,
		},
		{
			name: "Invalid email",
			args: args{
				email:    "invalid email",
				password: "test",
			},
			want:    user.User{},
			wantErr: user.ErrInvalidEmail,
		},
		{
			name: "Empty Password",
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
			got, err := user.NewUser(tt.args.email, tt.args.password)
			if err != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got.ID == uuid.Nil {
				t.Errorf("NewUser() got = %v, want %v", got, tt.want)
			}
			if tt.wantErr == nil && got.CreatedAt.IsZero() {
				t.Errorf("NewUser() CreatedAt is zero")
			}
			if tt.wantErr == nil && got.UpdatedAt.IsZero() {
				t.Errorf("NewUser() CreatedAt is UpdatedAt")
			}
			if tt.wantErr == nil && got.Email != tt.want.Email {
				t.Errorf("NewUser() got = %v, want %v", got.Email, tt.want.Email)
			}
			if tt.wantErr == nil && bcrypt.CompareHashAndPassword([]byte(got.Password), []byte(tt.args.password)) != nil {
				t.Error("NewUser() password is not encrypted correctly")
			}
		})
	}
}

func TestUserComparePassword(t *testing.T) {
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
				password: "qwerty123",
			},
			want: user.User{
				Email:    "test@example.com",
				Password: encryptString("qwerty123"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := user.NewUser(tt.args.email, tt.args.password)
			if err != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && (got.Password == tt.args.password || !got.ComparePassword(tt.args.password)) {
				t.Error("NewUser() password is not encrypted correctly")
			}
		})
	}
}

func encryptString(s string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}

	return string(b)
}
