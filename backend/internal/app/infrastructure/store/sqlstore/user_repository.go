package sqlstore

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/todolist"
)

type UserRepository struct {
	store *Store
}

type DBUserRecord struct {
	ID                uuid.UUID
	Email             string
	EncryptedPassword string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (r *DBUserRecord) ToUser() (*todolist.User, error) {
	u := &todolist.User{}

	if err := u.SetID(r.ID); err != nil {
		return nil, err
	}

	if err := u.SetEmail(r.Email); err != nil {
		return nil, err
	}

	if err := u.SetEncryptedPassword(r.EncryptedPassword); err != nil {
		return nil, err
	}

	if err := u.SetCreatedAt(r.CreatedAt); err != nil {
		return nil, err
	}

	if err := u.SetUpdatedAt(r.UpdatedAt); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) SaveUser(ctx context.Context, user *todolist.User) error {
	_, err := r.store.db.Exec(
		`
			INSERT INTO users (id, email, encrypted_password, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT(id) DO UPDATE
			SET email = EXCLUDED.email,
				encrypted_password = EXCLUDED.encrypted_password,
				updated_at = EXCLUDED.updated_at;
		`,
		user.GetID(),
		user.GetEmail(),
		user.GetEncryptedPassword(),
		user.GetCreatedAt(),
		user.GetUpdatedAt(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*todolist.User, error) {
	row := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password, created_at, updated_at FROM users WHERE email = $1;",
		email,
	)

	rec := &DBUserRecord{}
	err := row.Scan(
		&rec.ID,
		&rec.Email,
		&rec.EncryptedPassword,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return rec.ToUser()
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*todolist.User, error) {
	row := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password, created_at, updated_at FROM users WHERE id = $1;",
		id,
	)

	rec := &DBUserRecord{}
	err := row.Scan(
		&rec.ID,
		&rec.Email,
		&rec.EncryptedPassword,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return rec.ToUser()
}
