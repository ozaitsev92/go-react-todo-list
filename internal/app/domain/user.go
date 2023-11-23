package domain

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id                uuid.UUID
	email             string
	password          string
	encryptedPassword string
	createdAt         time.Time
	updatedAt         time.Time
}

func CreateUser(email, password string) (*User, error) {
	encryptedPassword, err := encryptString(password)
	if err != nil {
		return nil, err
	}

	currTime := time.Now()

	t := &User{
		id:                uuid.New(),
		email:             email,
		password:          password,
		encryptedPassword: encryptedPassword,
		createdAt:         currTime,
		updatedAt:         currTime,
	}

	if err := t.Validate(); err != nil {
		return nil, err
	}

	return t, nil
}

func (u *User) GetID() uuid.UUID {
	return u.id
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetCreatedAt() time.Time {
	return u.createdAt
}

func (u *User) GetUpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.email, validation.Required, is.Email),
		validation.Field(&u.password, validation.By(requiredIf(u.encryptedPassword == "")), validation.Length(6, 100)),
	)
}

func (u *User) BeforeUpdate() error {
	u.updatedAt = time.Now()

	return nil
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.encryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
