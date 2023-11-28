package todolist

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

func (u *User) SetID(id uuid.UUID) error {
	u.id = id
	return nil
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) SetEmail(email string) error {
	err := validation.Validate(email, validation.Required, is.Email)
	if err != nil {
		return err
	}

	u.email = email
	return nil
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) SetPassword(password string) error {
	err := validation.Validate(password, validation.By(requiredIf(u.encryptedPassword == "")), validation.Length(6, 100))
	if err != nil {
		return err
	}

	u.password = password
	return nil
}

func (u *User) GetEncryptedPassword() string {
	return u.encryptedPassword
}

func (u *User) SetEncryptedPassword(encryptedPassword string) error {
	err := validation.Validate(encryptedPassword, validation.Required)
	if err != nil {
		return err
	}

	u.encryptedPassword = encryptedPassword
	return nil
}

func (u *User) GetCreatedAt() time.Time {
	return u.createdAt
}

func (u *User) SetCreatedAt(createdAt time.Time) error {
	err := validation.Validate(createdAt, validation.By(timeNotZero))
	if err != nil {
		return err
	}

	u.createdAt = createdAt
	return nil
}

func (u *User) GetUpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) SetUpdatedAt(createdAt time.Time) error {
	err := validation.Validate(createdAt, validation.By(timeNotZero))
	if err != nil {
		return err
	}

	u.createdAt = createdAt
	return nil
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.email, validation.Required, is.Email),
		validation.Field(&u.password, validation.By(requiredIf(u.encryptedPassword == "")), validation.Length(6, 100)),
		validation.Field(&u.encryptedPassword, validation.Required),
		validation.Field(&u.createdAt, validation.By(timeNotZero)),
		validation.Field(&u.updatedAt, validation.By(timeNotZero)),
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
