package repository_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ozaitsev92/tododdd/config"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	repository "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/mongo"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/mongo/converter"
	repoModel "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/mongo/model"
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

	u, err := user.NewUser("test1@example.com", "Password123")
	if err != nil {
		t.Errorf("GetByID() failed to create a new user: err = '%v'", err)
	}

	// Check if a user exists in the DB: should fail
	r := repository.NewRepository(cfg)
	_, err = r.GetByID(context.Background(), u.ID)
	if !errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("GetByID() got = '%v', want = '%v'", err, user.ErrUserNotFound)
	}

	// Add the user to the users collection
	collection := mongodb.NewOrGetSingleton(cfg).Collection("users")
	mongoUser := converter.ToRepoFromUser(u)
	_, err = collection.InsertOne(context.Background(), mongoUser)
	if err != nil {
		t.Errorf("GetByID() failed to save a new user: err = '%v'", err)
	}

	// Check if a user exists in the DB: should succeed
	foundUser, err := r.GetByID(context.Background(), u.ID)
	if errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("GetByID() err = '%v', want = '%v'", err, nil)
	}

	if foundUser.ID != u.ID {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundUser.ID, u.ID)
	}

	if foundUser.Email != u.Email {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundUser.Email, u.Email)
	}
}

func TestRepositoryGetByEmail(t *testing.T) {
	cfg := config.Config{}
	cfg.MongoDBName = "todo_test"
	cfg.MongoUrl = fmt.Sprintf("mongodb://localhost:%s", MONGODB_PORT)

	u, err := user.NewUser("test2@example.com", "Password123")
	if err != nil {
		t.Errorf("GetByEmail()failed to create a new user: err = '%v'", err)
	}

	// Check if a user exists in the DB: should fail
	r := repository.NewRepository(cfg)
	_, err = r.GetByEmail(context.Background(), u.Email)
	if !errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("GetByEmail() got = '%v', want = '%v'", err, user.ErrUserNotFound)
	}

	// Add the user to the users collection
	collection := mongodb.NewOrGetSingleton(cfg).Collection("users")
	mongoUser := converter.ToRepoFromUser(u)
	_, err = collection.InsertOne(context.Background(), mongoUser)
	if err != nil {
		t.Errorf("GetByEmail()failed to save a new user: err = '%v'", err)
	}

	// Check if a user exists in the DB: should succeed
	foundUser, err := r.GetByEmail(context.Background(), u.Email)
	if errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("GetByID() err = '%v', want = '%v'", err, nil)
	}

	if foundUser.ID != u.ID {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundUser.ID, u.ID)
	}

	if foundUser.Email != u.Email {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundUser.Email, u.Email)
	}
}

func TestRepositorySave(t *testing.T) {
	cfg := config.Config{}
	cfg.MongoDBName = "todo_test"
	cfg.MongoUrl = fmt.Sprintf("mongodb://localhost:%s", MONGODB_PORT)

	u, err := user.NewUser("test3@example.com", "Password123")
	if err != nil {
		t.Errorf("Save() failed to create a new user: err = '%v'", err)
	}

	var mongoUser repoModel.User
	collection := mongodb.NewOrGetSingleton(cfg).Collection("users")

	// Check if a user exists in the DB: should fail
	err = collection.FindOne(context.Background(), bson.M{"email": u.Email}).Decode(&mongoUser)
	if !errors.Is(err, mongo.ErrNoDocuments) {
		t.Errorf("Save() got = '%v', want = '%v'", err, user.ErrUserNotFound)
	}

	// Save the user into the DB
	r := repository.NewRepository(cfg)
	err = r.Save(context.Background(), u)
	if err != nil {
		t.Errorf("Save() err = '%v'", err)
	}

	// Check if a user exists in the DB: should succeed
	err = collection.FindOne(context.Background(), bson.M{"email": u.Email}).Decode(&mongoUser)
	if err != nil {
		t.Errorf("Save() failed to create a new user: err = '%v'", err)
	}

	foundUser := converter.ToUserFromRepo(mongoUser)

	if foundUser.ID != u.ID {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundUser.ID, u.ID)
	}

	if foundUser.Email != u.Email {
		t.Errorf("GetByID() got = '%v', want = '%v'", foundUser.Email, u.Email)
	}
}
