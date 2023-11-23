package sqlstore

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/domain"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *domain.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	row := r.store.db.QueryRow(
		"INSERT INTO users (id, email, encrypted_password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		u.ID,
		u.Email,
		u.EncryptedPassword,
		u.CreatedAt,
		u.UpdatedAt,
	)

	return row.Scan(&u.ID)
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	u := &domain.User{}

	row := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE email = $1;",
		email,
	)

	if err := row.Scan(&u.ID, &u.Email, &u.EncryptedPassword); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) Find(id uuid.UUID) (*domain.User, error) {
	u := &domain.User{}

	row := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE id = $1;",
		id,
	)

	if err := row.Scan(&u.ID, &u.Email, &u.EncryptedPassword); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}
