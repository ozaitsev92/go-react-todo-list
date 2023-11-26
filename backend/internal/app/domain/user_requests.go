package domain

import "github.com/google/uuid"

type CreateUserRequest struct {
	Email    string
	Password string
}

type AuthenticateUserRequest struct {
	Email    string
	Password string
}

type FindUserByIDRequest struct {
	ID       uuid.UUID
	Password string
}
