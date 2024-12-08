package v1_test

import (
	"fmt"
	"log"
	"net"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ory/dockertest/v3"
	"github.com/ozaitsev92/tododdd/config"
	v1 "github.com/ozaitsev92/tododdd/internal/controller/http/v1"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/jwt"
	"github.com/ozaitsev92/tododdd/internal/usecase"

	taskRepository "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo"
	userRepository "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/mongo"
)

var (
	MONGODB_PORT = ""
)

type mockLogger struct{}

func (l *mockLogger) Debug(message interface{}, args ...interface{}) {}
func (l *mockLogger) Info(message string, args ...interface{})       {}
func (l *mockLogger) Warn(message string, args ...interface{})       {}
func (l *mockLogger) Error(message interface{}, args ...interface{}) {}
func (l *mockLogger) Fatal(message interface{}, args ...interface{}) {}

func setNewRouter() *gin.Engine {
	cfg := config.Config{}
	cfg.MongoDBName = "todo_test"
	cfg.MongoUrl = fmt.Sprintf("mongodb://localhost:%s", MONGODB_PORT)

	l := new(mockLogger)

	taskUseCase := usecase.NewTaskUseCase(
		taskRepository.NewRepository(cfg),
	)

	userUseCase := usecase.NewUserUseCase(
		userRepository.NewRepository(cfg),
	)

	jwtService := jwt.NewJWTService(
		[]byte(cfg.JWTSigningKey),
		cfg.JWTSessionLength,
		cfg.JWTCookieDomain,
		cfg.JWTSecureCookie,
	)

	gin.SetMode(gin.TestMode)

	handler := gin.Default()

	v1.NewRouter(handler, cfg, l, jwtService, taskUseCase, userUseCase)

	return handler
}

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
		_, err := net.Dial("tcp", net.JoinHostPort("localhost", MONGODB_PORT))
		return err
	})

	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestRepositoryHealthz(t *testing.T) {
	router := setNewRouter()

	req := httptest.NewRequest("GET", "/healthz", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("/healthz got = '%v', want = '%v'", w.Code, 200)
	}
}

func TestRepositoryMetrics(t *testing.T) {
	router := setNewRouter()

	req := httptest.NewRequest("GET", "/metrics", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("/metrics got = '%v', want = '%v'", w.Code, 200)
	}
}
