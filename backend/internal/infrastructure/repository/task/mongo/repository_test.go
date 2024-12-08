package repository_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ozaitsev92/tododdd/config"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	repository "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo/converter"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo/model"
	repoModel "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo/model"
	"github.com/ozaitsev92/tododdd/pkg/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var (
	MONGODB_PORT = ""
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("mongo", "latest", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	err = pool.Retry(func() error {
		MONGODB_PORT = resource.GetPort("27017/tcp")
		// ping to ensure that the server is up and running
		_, err := net.Dial("tcp", net.JoinHostPort("localhost", MONGODB_PORT))
		return err
	})

	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestRepositoryGetByID(t *testing.T) {
	cfg := config.Config{}
	cfg.MongoDBName = "todo_test"
	cfg.MongoUrl = fmt.Sprintf("mongodb://localhost:%s", MONGODB_PORT)

	ti, err := task.NewTask("task text", uuid.New())
	if err != nil {
		t.Errorf("GetByID() failed to create a new task: err = '%v'", err)
	}

	// Check if a task exists in the DB: should fail
	r := repository.NewRepository(cfg)
	_, err = r.GetByID(context.Background(), ti.ID)
	if !errors.Is(err, task.ErrTaskNotFound) {
		t.Errorf("GetByID() got = '%v', want = '%v'", err, task.ErrTaskNotFound)
	}

	// Add the task to the tasks collection
	collection := mongodb.NewOrGetSingleton(cfg).Collection("tasks")
	mongoTask := converter.ToRepoFromTask(ti)
	_, err = collection.InsertOne(context.Background(), mongoTask)
	if err != nil {
		t.Errorf("GetByID() failed to save a new task: err = '%v'", err)
	}

	// Check if a task exists in the DB: should succeed
	foundTask, err := r.GetByID(context.Background(), ti.ID)
	if errors.Is(err, task.ErrTaskNotFound) {
		t.Errorf("GetByID() err = '%v', want = '%v'", err, nil)
	}

	if foundTask.ID != ti.ID {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundTask.ID, ti.ID)
	}

	if foundTask.Text != ti.Text {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundTask.Text, ti.Text)
	}

	if foundTask.UserID != ti.UserID {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundTask.UserID, ti.UserID)
	}
}

func TestRepositoryGetAllByUserID(t *testing.T) {
	cfg := config.Config{}
	cfg.MongoDBName = "todo_test"
	cfg.MongoUrl = fmt.Sprintf("mongodb://localhost:%s", MONGODB_PORT)

	userId1 := uuid.New()

	t1, err := task.NewTask("task text 1", userId1)
	if err != nil {
		t.Errorf("GetAllByUserID() failed to create a new task: err = '%v'", err)
	}

	t2, err := task.NewTask("task text 2", userId1)
	if err != nil {
		t.Errorf("GetAllByUserID() failed to create a new task: err = '%v'", err)
	}

	userId2 := uuid.New()

	t3, err := task.NewTask("task text 2", userId2)
	if err != nil {
		t.Errorf("GetAllByUserID() failed to create a new task: err = '%v'", err)
	}

	// Check if a task exists in the DB: should fail
	r := repository.NewRepository(cfg)
	foundTasks, err := r.GetAllByUserID(context.Background(), t1.UserID)
	if err != nil {
		t.Errorf("GetAllByUserID() error = '%v'", err)
	}

	if len(foundTasks) != 0 {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", len(foundTasks), 2)
	}

	// Add the task to the tasks collection
	collection := mongodb.NewOrGetSingleton(cfg).Collection("tasks")
	_, err = collection.InsertMany(context.Background(), []interface{}{
		converter.ToRepoFromTask(t1),
		converter.ToRepoFromTask(t2),
		converter.ToRepoFromTask(t3),
	})
	if err != nil {
		t.Errorf("GetAllByUserID() failed to save new tasks: err = '%v'", err)
	}

	// Check if a task exists in the DB: should succeed
	foundTasks, err = r.GetAllByUserID(context.Background(), userId1)
	if errors.Is(err, task.ErrTaskNotFound) {
		t.Errorf("GetAllByUserID() err = '%v', want = '%v'", err, nil)
	}

	if len(foundTasks) != 2 {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", len(foundTasks), 2)
	}

	// Task #1
	ft1 := foundTasks[0]
	if ft1.ID != t1.ID {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", ft1.ID, t1.ID)
	}

	if ft1.Text != t1.Text {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", ft1.Text, t1.Text)
	}

	if ft1.UserID != t1.UserID {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", ft1.UserID, t1.UserID)
	}

	// Task #2
	ft2 := foundTasks[1]
	if ft2.ID != t2.ID {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", ft2.ID, t2.ID)
	}

	if ft2.Text != t2.Text {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", ft2.Text, t2.Text)
	}

	if ft2.UserID != t2.UserID {
		t.Errorf("GetAllByUserID() got = '%v', want = '%v'", ft2.UserID, t2.UserID)
	}
}

func TestRepositorySave(t *testing.T) {
	cfg := config.Config{}
	cfg.MongoDBName = "todo_test"
	cfg.MongoUrl = fmt.Sprintf("mongodb://localhost:%s", MONGODB_PORT)

	ti, err := task.NewTask("task text", uuid.New())
	if err != nil {
		t.Errorf("Save() failed to create a new task: err = '%v'", err)
	}

	var mongotask repoModel.Task
	collection := mongodb.NewOrGetSingleton(cfg).Collection("tasks")

	// Check if a task exists in the DB: should fail
	err = collection.FindOne(context.Background(), bson.M{"_id": ti.ID.String()}).Decode(&mongotask)
	if !errors.Is(err, mongo.ErrNoDocuments) {
		t.Errorf("Save() got = '%v', want = '%v'", err, user.ErrUserNotFound)
	}

	// Save the task into the DB
	r := repository.NewRepository(cfg)
	err = r.Save(context.Background(), ti)
	if err != nil {
		t.Errorf("Save() err = '%v'", err)
	}

	// Check if a user exists in the DB: should succeed
	err = collection.FindOne(context.Background(), bson.M{"_id": ti.ID.String()}).Decode(&mongotask)
	if err != nil {
		t.Errorf("Save() failed to create a new user: err = '%v'", err)
	}

	fountTask := converter.ToTaskFromRepo(mongotask)

	if fountTask.ID != ti.ID {
		t.Errorf("GetByID() got = '%v', want = '%v'", fountTask.ID, ti.ID)
	}

	if fountTask.UserID != ti.UserID {
		t.Errorf("GetByID() got = '%v', want = '%v'", fountTask.UserID, ti.UserID)
	}
}

func TestRepositoryUpdate(t *testing.T) {
	cfg := config.Config{}
	cfg.MongoDBName = "todo_test"
	cfg.MongoUrl = fmt.Sprintf("mongodb://localhost:%s", MONGODB_PORT)

	ti, err := task.NewTask("task text", uuid.New())
	if err != nil {
		t.Errorf("Update() failed to create a new task: err = '%v'", err)
	}

	// Add the task to the tasks collection
	collection := mongodb.NewOrGetSingleton(cfg).Collection("tasks")
	mongoTask := converter.ToRepoFromTask(ti)
	_, err = collection.InsertOne(context.Background(), mongoTask)
	if err != nil {
		t.Errorf("Update() failed to save a new task: err = '%v'", err)
	}

	ti.Text = "this is a new version of the task"

	// Update the task in the DB
	r := repository.NewRepository(cfg)
	err = r.Update(context.Background(), ti)
	if err != nil {
		t.Errorf("Update() err = '%v'", err)
	}

	var mongotask model.Task
	err = collection.FindOne(context.Background(), bson.M{"_id": ti.ID.String()}).Decode(&mongotask)
	if err != nil {
		t.Errorf("Save() failed to create a new user: err = '%v'", err)
	}

	fountTask := converter.ToTaskFromRepo(mongotask)

	if fountTask.ID != ti.ID {
		t.Errorf("GetByID() got = '%v', want = '%v'", fountTask.ID, ti.ID)
	}

	if fountTask.UserID != ti.UserID {
		t.Errorf("GetByID() got = '%v', want = '%v'", fountTask.UserID, ti.UserID)
	}

	if fountTask.Text != ti.Text {
		t.Errorf("GetByID() got = '%v', want = '%v'", fountTask.Text, ti.Text)
	}
}

func TestRepositoryDelete(t *testing.T) {
	cfg := config.Config{}
	cfg.MongoDBName = "todo_test"
	cfg.MongoUrl = fmt.Sprintf("mongodb://localhost:%s", MONGODB_PORT)

	ti, err := task.NewTask("task text", uuid.New())
	if err != nil {
		t.Errorf("Delete() failed to create a new task: err = '%v'", err)
	}

	// Add the task to the tasks collection
	collection := mongodb.NewOrGetSingleton(cfg).Collection("tasks")
	mongoTask := converter.ToRepoFromTask(ti)
	_, err = collection.InsertOne(context.Background(), mongoTask)
	if err != nil {
		t.Errorf("Delete() failed to save a new task: err = '%v'", err)
	}

	// Delete the task from the the DB
	r := repository.NewRepository(cfg)
	err = r.Delete(context.Background(), ti.ID)
	if err != nil {
		t.Errorf("Delete() err = '%v'", err)
	}

	// Check if a task exists in the DB: should fail
	var mongotask repoModel.Task
	err = collection.FindOne(context.Background(), bson.M{"_id": ti.ID.String()}).Decode(&mongotask)
	if !errors.Is(err, mongo.ErrNoDocuments) {
		t.Errorf("Save() got = '%v', want = '%v'", err, user.ErrUserNotFound)
	}
}
