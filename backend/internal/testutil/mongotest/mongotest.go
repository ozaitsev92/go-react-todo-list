//go:build integration

// Package mongotest provides MongoDB for integration tests: a dockertest
// container, or the server at MONGO_TEST_URI. Each test gets its own database.
package mongotest

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ozaitsev92/tododdd/pkg/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

// Image and Tag pin the test MongoDB; keep in sync with deployments.
const (
	Image = "mongo"
	Tag   = "8.0"
)

var (
	client *mongo.Client
	uri    string
)

// Main starts MongoDB and runs the tests; call it from TestMain.
func Main(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	uri = os.Getenv("MONGO_TEST_URI")

	if uri == "" {
		pool, err := dockertest.NewPool("")
		if err != nil {
			log.Printf("mongotest: docker unavailable: %v", err)
			return 1
		}
		pool.MaxWait = 60 * time.Second

		resource, err := pool.RunWithOptions(&dockertest.RunOptions{Repository: Image, Tag: Tag})
		if err != nil {
			log.Printf("mongotest: start container: %v", err)
			return 1
		}
		defer func() {
			if err := pool.Purge(resource); err != nil {
				log.Printf("mongotest: purge container: %v", err)
			}
		}()
		_ = resource.Expire(300) // in case the test binary is killed

		uri = "mongodb://" + resource.GetHostPort("27017/tcp")

		if err := pool.Retry(func() error {
			c, err := mongodb.Connect(context.Background(), uri, 2*time.Second)
			if err != nil {
				return err
			}
			client = c
			return nil
		}); err != nil {
			log.Printf("mongotest: mongo not ready: %v", err)
			return 1
		}
	} else {
		c, err := mongodb.Connect(context.Background(), uri, 10*time.Second)
		if err != nil {
			log.Printf("mongotest: connect %s: %v", uri, err)
			return 1
		}
		client = c
	}

	defer func() { _ = client.Disconnect(context.Background()) }()

	return m.Run()
}

// DB returns a fresh, uniquely named database that is dropped when t ends.
func DB(t testing.TB) *mongo.Database {
	t.Helper()

	if client == nil {
		t.Fatal("mongotest: client not initialised; call mongotest.Main from TestMain")
	}

	name := fmt.Sprintf("test_%s", strings.ReplaceAll(uuid.NewString(), "-", "")[:16])
	db := client.Database(name)

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := db.Drop(ctx); err != nil {
			t.Logf("mongotest: drop %s: %v", name, err)
		}
	})

	return db
}

// URI returns the test server's connection string.
func URI() string { return uri }
