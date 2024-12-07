package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserUseCase struct {
	userRepository user.Repository
}

// NewUserUseCase creates an new instance of the UserUseCase.
func NewUserUseCase(userRepository user.Repository) *UserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
	}
}

// RegisterNewUser registers a new user and saves it to the userRepository.
func (s *UserUseCase) RegisterNewUser(ctx context.Context, email, password string) (user.User, error) {
	_, err := s.userRepository.GetByEmail(ctx, email)

	if !errors.Is(err, user.ErrUserNotFound) {
		return user.User{}, ErrUserAlreadyExists
	}

	u, err := user.NewUser(email, password)

	if err != nil {
		return user.User{}, err
	}

	err = s.userRepository.Save(ctx, u)

	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

// GetUserByID returns a user by id.
func (s *UserUseCase) GetUserByID(ctx context.Context, id uuid.UUID) (user.User, error) {
	u, err := s.userRepository.GetByID(ctx, id)

	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

// GetUserByEmail returns a user by id.
func (s *UserUseCase) GetUserByEmail(ctx context.Context, email string) (user.User, error) {
	u, err := s.userRepository.GetByEmail(ctx, email)

	if err != nil {
		return user.User{}, err
	}

	return u, nil
}
