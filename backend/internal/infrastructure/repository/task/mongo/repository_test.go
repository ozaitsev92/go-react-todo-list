//go:build integration

package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/repository/contract"
	repository "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo"
	"github.com/ozaitsev92/tododdd/internal/testutil/mongotest"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMain(m *testing.M) { mongotest.Main(m) }

func newRepo(t *testing.T) *repository.Repository {
	t.Helper()
	r := repository.NewRepository(mongotest.DB(t))
	if err := r.EnsureIndexes(context.Background()); err != nil {
		t.Fatalf("EnsureIndexes: %v", err)
	}
	return r
}

func TestTaskRepository_Contract(t *testing.T) {
	contract.TaskRepository(t, func(t *testing.T) task.Repository { return newRepo(t) })
}

func TestTaskRepository_EnsureIndexes_CreatesUserIDIndex(t *testing.T) {
	db := mongotest.DB(t)
	r := repository.NewRepository(db)
	ctx := context.Background()

	for range 2 { // idempotent
		if err := r.EnsureIndexes(ctx); err != nil {
			t.Fatalf("EnsureIndexes: %v", err)
		}
	}

	cur, err := db.Collection(repository.CollectionName).Indexes().List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var idx []bson.M
	if err := cur.All(ctx, &idx); err != nil {
		t.Fatal(err)
	}
	for _, i := range idx {
		if i["name"] == "user_id_created_at" {
			return
		}
	}
	t.Fatalf("index user_id_created_at not found in %v", idx)
}

func TestTaskRepository_GetByID_CorruptDocument_ReturnsErrorNotPanic(t *testing.T) {
	db := mongotest.DB(t)
	r := repository.NewRepository(db)
	ctx := context.Background()
	id := uuid.New()

	_, err := db.Collection(repository.CollectionName).InsertOne(ctx,
		bson.M{"_id": id.String(), "text": "x", "user_id": "not-a-uuid"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = r.GetByID(ctx, id)
	if err == nil || errors.Is(err, task.ErrTaskNotFound) {
		t.Fatalf("err = %v, want a decode error", err)
	}
}
