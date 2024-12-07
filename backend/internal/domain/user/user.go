package user

import (
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail    = errors.New("email is invalid")
	ErrInvalidPassword = errors.New("password is invalid")
)

// User is a representation of a user entity.
type User struct {
	ID        uuid.UUID
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates and returns a new User.
func NewUser(email, password string) (User, error) {
	if email == "" {
		return User{}, ErrInvalidEmail
	}

	m, err := mail.ParseAddress(email)
	if err != nil {
		return User{}, ErrInvalidEmail
	}

	parsedEmail := m.Address

	if password == "" {
		return User{}, ErrInvalidPassword
	}

	encryptedPassword, err := encryptString(password)
	if err != nil {
		return User{}, err
	}

	currentTime := time.Now()

	user := User{
		ID:        uuid.New(),
		Email:     parsedEmail,
		Password:  encryptedPassword,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	return user, nil
}

// ComparePassword -.
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
