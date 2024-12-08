package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ory/dockertest/v3"
	"github.com/ozaitsev92/tododdd/config"
	v1 "github.com/ozaitsev92/tododdd/internal/controller/http/v1"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/jwt"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/model"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	"github.com/ozaitsev92/tododdd/internal/usecase"
	"github.com/ozaitsev92/tododdd/pkg/mongodb"

	taskRepository "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo"
	taskConverter "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo/converter"
	userRepository "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/mongo"
	userConverter "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/mongo/converter"
)

var (
	MONGODB_PORT = ""
)

const (
	jwtCookieName = "jwt-token"
)

type mockLogger struct{}

func (l *mockLogger) Debug(message interface{}, args ...interface{}) {}
func (l *mockLogger) Info(message string, args ...interface{})       {}
func (l *mockLogger) Warn(message string, args ...interface{})       {}
func (l *mockLogger) Error(message interface{}, args ...interface{}) {}
func (l *mockLogger) Fatal(message interface{}, args ...interface{}) {}

func newJsonRequest(method, url string, payload map[string]string) *http.Request {
	jsonPayload, _ := json.Marshal(payload)

	req := httptest.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func setNewRouter() (*gin.Engine, config.Config, *jwt.JWTService) {
	cfg := config.Config{}
	cfg.MongoDBName = "todo_test"
	cfg.MongoUrl = fmt.Sprintf("mongodb://localhost:%s", MONGODB_PORT)
	cfg.JWTSessionLength = 30

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

	return handler, cfg, jwtService
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
	router, _, _ := setNewRouter()

	req := newJsonRequest("GET", "/healthz", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("/healthz got = '%v', want = '%v'", w.Code, 200)
	}
}

func TestRepositoryMetrics(t *testing.T) {
	router, _, _ := setNewRouter()

	req := newJsonRequest("GET", "/metrics", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("/metrics got = '%v', want = '%v'", w.Code, 200)
	}
}

func TestRepositoryCreateUser(t *testing.T) {
	router, _, _ := setNewRouter()

	payload := map[string]string{
		"password": "Password123",
		"email":    "john@example.com",
	}
	req := newJsonRequest("POST", "/v1/users", payload)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("/v1/users got = '%v', want = '%v'", w.Code, 201)
	}

	var response model.User

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("/v1/users error = '%v'", err)
	}

	if response.Email != payload["email"] {
		t.Errorf("/v1/users got = '%v', want = '%v'", response.Email, payload["email"])
	}

	if response.ID == "" {
		t.Error("/v1/users User.ID is empty")
	}
}

func TestRepositoryLogin(t *testing.T) {
	router, cfg, _ := setNewRouter()

	rawPassword := "Password123"
	u, err := user.NewUser("test1@example.com", rawPassword)
	if err != nil {
		t.Errorf("/v1/users/login failed to create a new user: err = '%v'", err)
	}

	// Add the user to the users collection
	collection := mongodb.NewOrGetSingleton(cfg).Collection("users")
	mongoUser := userConverter.ToRepoFromUser(u)
	_, err = collection.InsertOne(context.Background(), mongoUser)
	if err != nil {
		t.Errorf("/v1/users/login failed to save a new user: err = '%v'", err)
	}

	payload := map[string]string{
		"password": rawPassword,
		"email":    u.Email,
	}
	req := newJsonRequest("POST", "/v1/users/login", payload)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("/v1/users/login got = '%v', want = '%v'", w.Code, 200)
	}

	if len(w.Body.Bytes()) == 0 {
		t.Errorf("/v1/users/login response body should be empty, got = '%s'", w.Body.String())
	}

	resp := w.Result()
	cookies := resp.Cookies()

	var testCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == jwtCookieName {
			testCookie = cookie
			break
		}
	}

	if testCookie == nil {
		t.Errorf("/v1/users/login response must have a '%s' cookie", jwtCookieName)
	}

	if testCookie != nil && testCookie.Value == "" {
		t.Errorf("/v1/users/login '%s' cookie must have a Value, got = '%s'", jwtCookieName, testCookie.Value)
	}

	if testCookie != nil && testCookie.MaxAge <= 0 {
		t.Errorf("/v1/users/login '%s' cookie must have a MaxAge > 0, got = '%d'", jwtCookieName, testCookie.MaxAge)
	}
}

func TestRepositoryLogout(t *testing.T) {
	router, cfg, _ := setNewRouter()

	u, err := user.NewUser("test1@example.com", "Password123")
	if err != nil {
		t.Errorf("/v1/users/logout failed to create a new user: err = '%v'", err)
	}

	// Add the user to the users collection
	collection := mongodb.NewOrGetSingleton(cfg).Collection("users")
	mongoUser := userConverter.ToRepoFromUser(u)
	_, err = collection.InsertOne(context.Background(), mongoUser)
	if err != nil {
		t.Errorf("/v1/users/logout failed to save a new user: err = '%v'", err)
	}

	req := newJsonRequest("POST", "/v1/users/logout", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("/v1/users/logout got = '%v', want = '%v'", w.Code, 200)
	}

	if len(w.Body.Bytes()) == 0 {
		t.Errorf("/v1/users/logout response body should be empty, got = '%s'", w.Body.String())
	}

	resp := w.Result()
	cookies := resp.Cookies()

	var testCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == jwtCookieName {
			testCookie = cookie
			break
		}
	}

	if testCookie == nil {
		t.Errorf("/v1/users/logout response must have a '%s' cookie", jwtCookieName)
	}

	if testCookie != nil && testCookie.Value != "" {
		t.Errorf("/v1/users/logout '%s' cookie must not have a Value, got = '%s'", jwtCookieName, testCookie.Value)
	}

	if testCookie != nil && testCookie.MaxAge != -1 {
		t.Errorf("/v1/users/logout '%s' cookie must have a MaxAge == -1, got = '%d'", jwtCookieName, testCookie.MaxAge)
	}
}

func TestRepositoryCurrent(t *testing.T) {
	router, cfg, jwtService := setNewRouter()

	u, err := user.NewUser("test1@example.com", "Password123")
	if err != nil {
		t.Errorf("/v1/users/current failed to create a new user: err = '%v'", err)
	}

	// Add the user to the users collection
	collection := mongodb.NewOrGetSingleton(cfg).Collection("users")
	mongoUser := userConverter.ToRepoFromUser(u)
	_, err = collection.InsertOne(context.Background(), mongoUser)
	if err != nil {
		t.Errorf("/v1/users/current failed to save a new user: err = '%v'", err)
	}

	token, err := jwtService.CreateJWTTokenForUser(u.ID)
	if err != nil {
		return
	}

	jwtCookie := jwtService.AuthCookie(token)

	req := newJsonRequest("GET", "/v1/users/current", nil)
	req.AddCookie(&http.Cookie{Name: jwtCookie.Name, Value: jwtCookie.Value})

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("/v1/users/current got = '%v', want = '%v'", w.Code, 200)
	}

	var response model.User

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("/v1/users/current error = '%v'", err)
	}

	if response.ID != u.ID.String() {
		t.Errorf("/v1/users/current got = '%v', want = '%v'", response.ID, u.ID)
	}

	if response.Email != u.Email {
		t.Errorf("/v1/users/current got = '%v', want = '%v'", response.Email, u.Email)
	}
}

func TestRepositoryCreateTask(t *testing.T) {
	router, cfg, jwtService := setNewRouter()

	u, err := user.NewUser("test1@example.com", "Password123")
	if err != nil {
		t.Errorf("/v1/tasks failed to create a new user: err = '%v'", err)
	}

	// Add the user to the users collection
	collection := mongodb.NewOrGetSingleton(cfg).Collection("users")
	mongoUser := userConverter.ToRepoFromUser(u)
	_, err = collection.InsertOne(context.Background(), mongoUser)
	if err != nil {
		t.Errorf("/v1/tasks failed to save a new user: err = '%v'", err)
	}

	token, err := jwtService.CreateJWTTokenForUser(u.ID)
	if err != nil {
		return
	}

	jwtCookie := jwtService.AuthCookie(token)

	payload := map[string]string{
		"text": "task text",
	}
	req := newJsonRequest("POST", "/v1/tasks", payload)
	req.AddCookie(&http.Cookie{Name: jwtCookie.Name, Value: jwtCookie.Value})

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("/v1/tasks got = '%v', want = '%v'", w.Code, 200)
	}

	var response model.Task

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("/v1/tasks error = '%v'", err)
	}

	if response.UserID != u.ID.String() {
		t.Errorf("/v1/tasks got = '%v', want = '%v'", response.ID, u.ID)
	}

	if response.Text != payload["text"] {
		t.Errorf("/v1/tasks got = '%v', want = '%v'", response.Text, payload["text"])
	}

	if response.Completed {
		t.Error("/v1/tasks task should not be completed")
	}
}

func TestRepositoryUpdateTask(t *testing.T) {
	router, cfg, jwtService := setNewRouter()

	// Add the user to the users collection
	u, err := user.NewUser("test1@example.com", "Password123")
	if err != nil {
		t.Errorf("/v1/tasks/:id failed to create a new user: err = '%v'", err)
	}

	mongoUser := userConverter.ToRepoFromUser(u)
	_, err = mongodb.NewOrGetSingleton(cfg).Collection("users").InsertOne(context.Background(), mongoUser)
	if err != nil {
		t.Errorf("/v1/tasks/:id failed to save a new user: err = '%v'", err)
	}

	// Add the task to the users collection
	ti, err := task.NewTask("task text v1", u.ID)
	if err != nil {
		t.Errorf("/v1/tasks/:id failed to create a new user: err = '%v'", err)
	}

	mongoTask := taskConverter.ToRepoFromTask(ti)
	_, err = mongodb.NewOrGetSingleton(cfg).Collection("tasks").InsertOne(context.Background(), mongoTask)
	if err != nil {
		t.Errorf("/v1/tasks/:id failed to save a new task: err = '%v'", err)
	}

	token, err := jwtService.CreateJWTTokenForUser(u.ID)
	if err != nil {
		return
	}

	jwtCookie := jwtService.AuthCookie(token)

	payload := map[string]string{
		"text": "task text v2",
	}
	req := newJsonRequest("PUT", "/v1/tasks/"+ti.ID.String(), payload)
	req.AddCookie(&http.Cookie{Name: jwtCookie.Name, Value: jwtCookie.Value})

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("/v1/tasks/:id got = '%v', want = '%v'", w.Code, 200)
	}

	var response model.Task

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("/v1/tasks/:id error = '%v'", err)
	}

	if response.UserID != u.ID.String() {
		t.Errorf("/v1/tasks/:id got = '%v', want = '%v'", response.ID, u.ID)
	}

	if response.Text != payload["text"] {
		t.Errorf("/v1/tasks/:id got = '%v', want = '%v'", response.Text, payload["text"])
	}

	if response.Completed {
		t.Error("/v1/tasks/:id task should not be completed")
	}
}

func TestRepositoryMarkTaskCompleted(t *testing.T) {
	router, cfg, jwtService := setNewRouter()

	// Add the user to the users collection
	u, err := user.NewUser("test1@example.com", "Password123")
	if err != nil {
		t.Errorf("/v1/tasks/:id/mark-completed failed to create a new user: err = '%v'", err)
	}

	mongoUser := userConverter.ToRepoFromUser(u)
	_, err = mongodb.NewOrGetSingleton(cfg).Collection("users").InsertOne(context.Background(), mongoUser)
	if err != nil {
		t.Errorf("/v1/tasks/:id/mark-completed failed to save a new user: err = '%v'", err)
	}

	// Add the task to the users collection
	ti, err := task.NewTask("task text v1", u.ID)
	if err != nil {
		t.Errorf("/v1/tasks/:id/mark-completed failed to create a new user: err = '%v'", err)
	}

	mongoTask := taskConverter.ToRepoFromTask(ti)
	_, err = mongodb.NewOrGetSingleton(cfg).Collection("tasks").InsertOne(context.Background(), mongoTask)
	if err != nil {
		t.Errorf("/v1/tasks/:id/mark-completed failed to save a new task: err = '%v'", err)
	}

	token, err := jwtService.CreateJWTTokenForUser(u.ID)
	if err != nil {
		return
	}

	jwtCookie := jwtService.AuthCookie(token)

	req := newJsonRequest("PUT", "/v1/tasks/"+ti.ID.String()+"/mark-completed", nil)
	req.AddCookie(&http.Cookie{Name: jwtCookie.Name, Value: jwtCookie.Value})

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("/v1/tasks/:id/mark-completed got = '%v', want = '%v'", w.Code, 200)
	}

	var response model.Task

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("/v1/tasks/:id/mark-completed error = '%v'", err)
	}

	if response.UserID != u.ID.String() {
		t.Errorf("/v1/tasks/:id/mark-completed got = '%v', want = '%v'", response.ID, u.ID)
	}

	if response.Text != ti.Text {
		t.Errorf("/v1/tasks/:id/mark-completed got = '%v', want = '%v'", response.Text, ti.Text)
	}

	if !response.Completed {
		t.Error("/v1/tasks/:id/mark-completed task should be completed")
	}
}

func TestRepositoryMarkTaskNotCompleted(t *testing.T) {
	router, cfg, jwtService := setNewRouter()

	// Add the user to the users collection
	u, err := user.NewUser("test1@example.com", "Password123")
	if err != nil {
		t.Errorf("/v1/tasks/:id/mark-not-completed failed to create a new user: err = '%v'", err)
	}

	mongoUser := userConverter.ToRepoFromUser(u)
	_, err = mongodb.NewOrGetSingleton(cfg).Collection("users").InsertOne(context.Background(), mongoUser)
	if err != nil {
		t.Errorf("/v1/tasks/:id/mark-not-completed failed to save a new user: err = '%v'", err)
	}

	// Add the task to the users collection
	ti, err := task.NewTask("task text v1", u.ID)
	if err != nil {
		t.Errorf("/v1/tasks/:id/mark-not-completed failed to create a new user: err = '%v'", err)
	}

	ti.Completed = true
	mongoTask := taskConverter.ToRepoFromTask(ti)
	_, err = mongodb.NewOrGetSingleton(cfg).Collection("tasks").InsertOne(context.Background(), mongoTask)
	if err != nil {
		t.Errorf("/v1/tasks/:id/mark-not-completed failed to save a new task: err = '%v'", err)
	}

	token, err := jwtService.CreateJWTTokenForUser(u.ID)
	if err != nil {
		return
	}

	jwtCookie := jwtService.AuthCookie(token)

	req := newJsonRequest("PUT", "/v1/tasks/"+ti.ID.String()+"/mark-not-completed", nil)
	req.AddCookie(&http.Cookie{Name: jwtCookie.Name, Value: jwtCookie.Value})

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("/v1/tasks/:id/mark-not-completed got = '%v', want = '%v'", w.Code, 200)
	}

	var response model.Task

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("/v1/tasks/:id/mark-not-completed error = '%v'", err)
	}

	if response.UserID != u.ID.String() {
		t.Errorf("/v1/tasks/:id/mark-not-completed got = '%v', want = '%v'", response.ID, u.ID)
	}

	if response.Text != ti.Text {
		t.Errorf("/v1/tasks/:id/mark-not-completed got = '%v', want = '%v'", response.Text, ti.Text)
	}

	if response.Completed {
		t.Error("/v1/tasks/:id/mark-not-completed task should not be completed")
	}
}

func TestRepositoryDeleteTask(t *testing.T) {
	router, cfg, jwtService := setNewRouter()

	// Add the user to the users collection
	u, err := user.NewUser("test1@example.com", "Password123")
	if err != nil {
		t.Errorf("/v1/tasks/:id failed to create a new user: err = '%v'", err)
	}

	mongoUser := userConverter.ToRepoFromUser(u)
	_, err = mongodb.NewOrGetSingleton(cfg).Collection("users").InsertOne(context.Background(), mongoUser)
	if err != nil {
		t.Errorf("/v1/tasks/:id failed to save a new user: err = '%v'", err)
	}

	// Add the task to the users collection
	ti, err := task.NewTask("task text v1", u.ID)
	if err != nil {
		t.Errorf("/v1/tasks/:id failed to create a new user: err = '%v'", err)
	}

	ti.Completed = true
	mongoTask := taskConverter.ToRepoFromTask(ti)
	_, err = mongodb.NewOrGetSingleton(cfg).Collection("tasks").InsertOne(context.Background(), mongoTask)
	if err != nil {
		t.Errorf("/v1/tasks/:id failed to save a new task: err = '%v'", err)
	}

	token, err := jwtService.CreateJWTTokenForUser(u.ID)
	if err != nil {
		return
	}

	jwtCookie := jwtService.AuthCookie(token)

	req := newJsonRequest("DELETE", "/v1/tasks/"+ti.ID.String(), nil)
	req.AddCookie(&http.Cookie{Name: jwtCookie.Name, Value: jwtCookie.Value})

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("/v1/tasks/:id got = '%v', want = '%v'", w.Code, 200)
	}

	if len(w.Body.Bytes()) == 0 {
		t.Errorf("/v1/tasks/:id response body should be empty, got = '%s'", w.Body.String())
	}
}
