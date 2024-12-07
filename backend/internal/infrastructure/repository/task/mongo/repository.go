package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/config"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo/converter"
	"github.com/ozaitsev92/tododdd/pkg/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"

	repoModel "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo/model"
)

var _ task.Repository = (*Repository)(nil)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(cfg config.Config) *Repository {
	collection := mongodb.NewOrGetSingleton(cfg).Collection("tasks")

	return &Repository{
		collection: collection,
	}
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (task.Task, error) {
	var mongoTask repoModel.Task

	err := r.collection.FindOne(ctx, bson.M{"_id": id.String()}).Decode(&mongoTask)
	if err != nil {
		return converter.ToTaskFromRepo(mongoTask), err
	}

	return converter.ToTaskFromRepo(mongoTask), nil
}

func (r *Repository) GetAllByUserID(ctx context.Context, userId uuid.UUID) ([]task.Task, error) {
	filter := bson.M{"user_id": userId.String()}
	sort := options.Find().SetSort(bson.M{"created_at": 1})

	cursor, err := r.collection.Find(ctx, filter, sort)
	if err != nil {
		return []task.Task{}, err
	}

	var mongoTasks []repoModel.Task

	err = cursor.All(ctx, &mongoTasks)
	if err != nil {
		return []task.Task{}, err
	}

	tasks := make([]task.Task, 0, len(mongoTasks))
	for _, mongoTask := range mongoTasks {
		tasks = append(tasks, converter.ToTaskFromRepo(mongoTask))
	}

	return tasks, nil
}

func (r *Repository) Save(ctx context.Context, t task.Task) error {
	mongoItem := converter.ToRepoFromTask(t)
	_, err := r.collection.InsertOne(ctx, mongoItem)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, t task.Task) error {
	filter := bson.M{"_id": t.ID.String()}
	update := bson.M{
		"$set": bson.M{
			"text":      t.Text,
			"completed": t.Completed,
		},
	}

	result := r.collection.FindOneAndUpdate(ctx, filter, update)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return task.ErrTaskNotFound
		}

		return task.ErrFailedUpdateTask
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.collection.FindOneAndDelete(ctx, bson.M{"_id": id.String()})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return task.ErrTaskNotFound
		}

		return task.ErrFailedDeleteTask
	}

	return nil
}
