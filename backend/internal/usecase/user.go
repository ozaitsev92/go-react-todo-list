package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
)

// ErrInvalidCredentials is returned for an unknown email or a wrong password.
var ErrInvalidCredentials = errors.New("invalid email or password")

// PasswordHasher hashes passwords and verifies them in constant time.
type PasswordHasher interface {
	Hash(password string) (string, error)
	// Compare returns nil if password matches hash.
	Compare(hash, password string) error
}

// UserUseCase registers and authenticates users.
type UserUseCase struct {
	users     user.Repository
	hasher    PasswordHasher
	dummyHash string
}

// NewUserUseCase creates a UserUseCase. The dummy hash makes failed logins for
// unknown emails as slow as wrong passwords.
func NewUserUseCase(users user.Repository, hasher PasswordHasher) (*UserUseCase, error) {
	dummy, err := hasher.Hash("dummy-password-for-timing")
	if err != nil {
		return nil, fmt.Errorf("new user use case: %w", err)
	}

	return &UserUseCase{users: users, hasher: hasher, dummyHash: dummy}, nil
}

// RegisterNewUser validates and stores a new user.
func (s *UserUseCase) RegisterNewUser(ctx context.Context, email, password string) (user.User, error) {
	email, err := user.NormalizeEmail(email)
	if err != nil {
		return user.User{}, err
	}

	if err := user.ValidatePassword(password); err != nil {
		return user.User{}, err
	}

	// Friendly fast path; the unique index is what guarantees uniqueness.
	_, err = s.users.GetByEmail(ctx, email)
	switch {
	case err == nil:
		return user.User{}, user.ErrEmailTaken
	case !errors.Is(err, user.ErrUserNotFound):
		return user.User{}, fmt.Errorf("register: look up email: %w", err)
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return user.User{}, fmt.Errorf("register: hash password: %w", err)
	}

	u, err := user.NewUser(email, hash)
	if err != nil {
		return user.User{}, err
	}

	if err := s.users.Save(ctx, u); err != nil {
		return user.User{}, fmt.Errorf("register: save user: %w", err)
	}

	return u, nil
}

// Authenticate returns the user matching the credentials.
func (s *UserUseCase) Authenticate(ctx context.Context, email, password string) (user.User, error) {
	normalized, err := user.NormalizeEmail(email)
	if err != nil {
		_ = s.hasher.Compare(s.dummyHash, password)
		return user.User{}, ErrInvalidCredentials
	}

	u, err := s.users.GetByEmail(ctx, normalized)
	if errors.Is(err, user.ErrUserNotFound) {
		_ = s.hasher.Compare(s.dummyHash, password)
		return user.User{}, ErrInvalidCredentials
	}
	if err != nil {
		return user.User{}, fmt.Errorf("authenticate: look up email: %w", err)
	}

	if err := s.hasher.Compare(u.PasswordHash, password); err != nil {
		return user.User{}, ErrInvalidCredentials
	}

	return u, nil
}

// GetUserByID returns the user with the given ID.
func (s *UserUseCase) GetUserByID(ctx context.Context, id uuid.UUID) (user.User, error) {
	u, err := s.users.GetByID(ctx, id)
	if err != nil {
		return user.User{}, fmt.Errorf("get user %s: %w", id, err)
	}

	return u, nil
}
