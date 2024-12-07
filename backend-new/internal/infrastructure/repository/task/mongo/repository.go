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

func (r *Repository) GetAllByUserID(ctx context.Context, UserID uuid.UUID) ([]task.Task, error) {
	filter := bson.M{"todo_list_id": UserID}
	sort := bson.D{{Name: "created_at", Value: 1}}
	opts := options.Find().SetSort(sort)

	cursor, err := r.collection.Find(ctx, filter, opts)
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
	result := r.collection.FindOneAndUpdate(ctx, bson.M{"_id": t.ID.String()}, converter.ToRepoFromTask(t))
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
