package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/config"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/mongo/converter"
	repoModel "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/mongo/model"
	"github.com/ozaitsev92/tododdd/pkg/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var _ user.Repository = (*Repository)(nil)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(cfg config.Config) *Repository {
	collection := mongodb.NewOrGetSingleton(cfg).Collection("users")

	return &Repository{
		collection: collection,
	}
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (user.User, error) {
	var u repoModel.User

	err := r.collection.FindOne(ctx, bson.M{"_id": id.String()}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.User{}, user.ErrUserNotFound
		}

		return user.User{}, err
	}

	return converter.ToUserFromRepo(u), nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (user.User, error) {
	var u repoModel.User

	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.User{}, user.ErrUserNotFound
		}

		return user.User{}, err
	}

	return converter.ToUserFromRepo(u), nil
}

func (r *Repository) Save(ctx context.Context, t user.User) error {
	mongoItem := converter.ToRepoFromUser(t)
	_, err := r.collection.InsertOne(ctx, mongoItem)
	if err != nil {
		return user.ErrFailedToStoreUser
	}

	return nil
}
