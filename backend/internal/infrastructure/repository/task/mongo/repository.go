// Package repository is the MongoDB task.Repository.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo/converter"
	repoModel "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CollectionName is the MongoDB collection holding tasks.
const CollectionName = "tasks"

var _ task.Repository = (*Repository)(nil)

// Repository stores tasks in a MongoDB collection.
type Repository struct {
	collection *mongo.Collection
}

// NewRepository returns a repository using db's tasks collection.
func NewRepository(db *mongo.Database) *Repository {
	return &Repository{collection: db.Collection(CollectionName)}
}

// EnsureIndexes creates the query indexes (idempotent).
func (r *Repository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: 1}},
		Options: options.Index().SetName("user_id_created_at"),
	})
	if err != nil {
		return fmt.Errorf("tasks: create indexes: %w", err)
	}

	return nil
}

// GetByID returns the task with id or task.ErrTaskNotFound.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (task.Task, error) {
	var doc repoModel.Task

	err := r.collection.FindOne(ctx, bson.M{"_id": id.String()}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return task.Task{}, task.ErrTaskNotFound
	}
	if err != nil {
		return task.Task{}, fmt.Errorf("tasks: find %s: %w", id, err)
	}

	return converter.ToTaskFromRepo(doc)
}

// GetAllByUserID returns the user's tasks ordered by created_at (never nil).
func (r *Repository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]task.Task, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID.String()}, opts)
	if err != nil {
		return nil, fmt.Errorf("tasks: find by user %s: %w", userID, err)
	}

	var docs []repoModel.Task
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("tasks: decode by user %s: %w", userID, err)
	}

	tasks := make([]task.Task, 0, len(docs))
	for _, doc := range docs {
		t, err := converter.ToTaskFromRepo(doc)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// Save inserts a new task or returns task.ErrTaskExists.
func (r *Repository) Save(ctx context.Context, t task.Task) error {
	_, err := r.collection.InsertOne(ctx, converter.ToRepoFromTask(t))
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("tasks: insert %s: %w", t.ID, task.ErrTaskExists)
	}
	if err != nil {
		return fmt.Errorf("tasks: insert %s: %w", t.ID, err)
	}

	return nil
}

// Update saves text, completed and updated_at, or returns task.ErrTaskNotFound.
func (r *Repository) Update(ctx context.Context, t task.Task) error {
	res, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": t.ID.String()},
		bson.M{"$set": bson.M{
			"text":       t.Text,
			"completed":  t.Completed,
			"updated_at": t.UpdatedAt,
		}},
	)
	if err != nil {
		return fmt.Errorf("tasks: update %s: %w", t.ID, err)
	}
	if res.MatchedCount == 0 {
		return task.ErrTaskNotFound
	}

	return nil
}

// Delete removes a task or returns task.ErrTaskNotFound.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": id.String()})
	if err != nil {
		return fmt.Errorf("tasks: delete %s: %w", id, err)
	}
	if res.DeletedCount == 0 {
		return task.ErrTaskNotFound
	}

	return nil
}
