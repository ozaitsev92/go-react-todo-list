package todolist

import (
	"context"
	"errors"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
)

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (ts *UserService) CreateUser(ctx context.Context, r *CreateUserRequest) (*User, error) {
	u, err := CreateUser(r.Email, r.Password)
	if err != nil {
		return nil, err
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := ts.userRepository.SaveUser(ctx, u); err != nil {
		return nil, err
	}

	return u, err
}

func (ts *UserService) AuthenticateUser(ctx context.Context, r *AuthenticateUserRequest) (*User, error) {
	u, err := ts.userRepository.FindByEmail(ctx, r.Email)
	if err != nil {
		return nil, errIncorrectEmailOrPassword
	}

	if !u.ComparePassword(r.Password) {
		return nil, errIncorrectEmailOrPassword
	}

	return u, nil
}

func (ts *UserService) FindUserByID(ctx context.Context, r *FindUserByIDRequest) (*User, error) {
	u, err := ts.userRepository.FindByID(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	return u, nil
}
