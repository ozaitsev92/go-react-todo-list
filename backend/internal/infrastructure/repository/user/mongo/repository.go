// Package repository is the MongoDB user.Repository.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/mongo/converter"
	repoModel "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/mongo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CollectionName is the MongoDB collection holding users.
const CollectionName = "users"

var _ user.Repository = (*Repository)(nil)

// Repository stores users in a MongoDB collection.
type Repository struct {
	collection *mongo.Collection
}

// NewRepository returns a repository using db's users collection.
func NewRepository(db *mongo.Database) *Repository {
	return &Repository{collection: db.Collection(CollectionName)}
}

// EnsureIndexes creates the unique email index (idempotent).
func (r *Repository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetName("email_unique").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("users: create indexes: %w", err)
	}

	return nil
}

// GetByID returns the user with id or user.ErrUserNotFound.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (user.User, error) {
	return r.findOne(ctx, bson.M{"_id": id.String()})
}

// GetByEmail returns the user with the (normalised) email or user.ErrUserNotFound.
func (r *Repository) GetByEmail(ctx context.Context, email string) (user.User, error) {
	return r.findOne(ctx, bson.M{"email": email})
}

// Save inserts a user or returns user.ErrEmailTaken on a duplicate email.
func (r *Repository) Save(ctx context.Context, u user.User) error {
	_, err := r.collection.InsertOne(ctx, converter.ToRepoFromUser(u))
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("users: insert: %w", user.ErrEmailTaken)
	}
	if err != nil {
		return fmt.Errorf("users: insert: %w", err)
	}

	return nil
}

func (r *Repository) findOne(ctx context.Context, filter bson.M) (user.User, error) {
	var doc repoModel.User

	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user.User{}, user.ErrUserNotFound
	}
	if err != nil {
		return user.User{}, fmt.Errorf("users: find: %w", err)
	}

	return converter.ToUserFromRepo(doc)
}
