//go:build integration

package app

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ozaitsev92/tododdd/config"
	v1 "github.com/ozaitsev92/tododdd/internal/controller/http/v1"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/jwt"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/password"
	"github.com/ozaitsev92/tododdd/internal/testutil/mongotest"
	"github.com/ozaitsev92/tododdd/internal/usecase"
	"github.com/ozaitsev92/tododdd/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) { mongotest.Main(m) }

// TestE2E_SignupLoginCRUD_AgainstMongo drives the real router over HTTP
// against MongoDB, including index creation via openStorage.
func TestE2E_SignupLoginCRUD_AgainstMongo(t *testing.T) {
	cfg := config.Default()
	log := logger.Discard()
	store := storageFor(t, mongotest.DB(t).Name())

	hasher, _ := password.NewBcryptHasher(bcrypt.MinCost)
	users, err := usecase.NewUserUseCase(store.users, hasher)
	if err != nil {
		t.Fatal(err)
	}
	tokens, err := jwt.NewService(jwt.Config{SigningKey: []byte(strings.Repeat("k", 32)), Issuer: "todo-api", TTL: time.Minute})
	if err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	srv := httptest.NewServer(v1.NewRouter(v1.Deps{
		Log: log, Tokens: tokens, Users: users, Tasks: usecase.NewTaskUseCase(store.tasks),
		MaxBodyBytes: cfg.MaxBodyBytes, Ready: store.ready,
	}))
	defer srv.Close()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	call := func(method, path, body string, want int) []byte {
		t.Helper()
		req, _ := http.NewRequestWithContext(t.Context(), method, srv.URL+path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != want {
			t.Fatalf("%s %s = %d, want %d (%s)", method, path, resp.StatusCode, want, b)
		}
		return b
	}

	creds := `{"email":"e2e@example.com","password":"Secret123"}`
	call(http.MethodGet, "/readyz", "", http.StatusOK)
	call(http.MethodPost, "/v1/users", creds, http.StatusCreated)
	call(http.MethodPost, "/v1/users", `{"email":"E2E@example.com","password":"Secret123"}`, http.StatusConflict)
	call(http.MethodPost, "/v1/users/login", creds, http.StatusOK)

	var created struct{ ID string }
	if err := json.Unmarshal(call(http.MethodPost, "/v1/tasks", `{"text":"e2e"}`, http.StatusCreated), &created); err != nil {
		t.Fatal(err)
	}
	call(http.MethodPut, "/v1/tasks/"+created.ID+"/mark-completed", "", http.StatusOK)

	var tasks []struct {
		ID        string
		Completed bool
	}
	if err := json.Unmarshal(call(http.MethodGet, "/v1/tasks", "", http.StatusOK), &tasks); err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 || tasks[0].ID != created.ID || !tasks[0].Completed {
		t.Fatalf("tasks = %+v", tasks)
	}

	call(http.MethodDelete, "/v1/tasks/"+created.ID, "", http.StatusNoContent)
	call(http.MethodPost, "/v1/users/logout", "", http.StatusNoContent)
	call(http.MethodGet, "/v1/tasks", "", http.StatusUnauthorized)
}

// storageFor runs openStorage against the shared test server and database.
func storageFor(t *testing.T, dbName string) storage {
	t.Helper()
	cfg := config.Default()
	cfg.MongoURL = mongotest.URI()
	cfg.MongoDBName = dbName

	s, err := openStorage(t.Context(), cfg, logger.Discard())
	if err != nil {
		t.Fatalf("openStorage: %v", err)
	}
	t.Cleanup(func() { s.close(context.Background(), logger.Discard()) })
	return s
}
