package apiserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/apiserver/jwt"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/domain"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestServer_HandleUsersCreate(t *testing.T) {
	jwtService := jwt.NewJWTService([]byte("test"), 30, "localhost", true)
	s := newServer(teststore.New(), jwtService)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "user@example.org",
				"password": "password",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			b := &bytes.Buffer{}
			assert.NoError(t, json.NewEncoder(b).Encode(tc.payload))

			req, _ := http.NewRequest(http.MethodPost, "/users", b)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleUserLogin(t *testing.T) {
	u := domain.TestUser(t, "email@example.com", "a password")

	store := teststore.New()
	assert.NoError(t, store.User().SaveUser(context.Background(), u))

	jwtService := jwt.NewJWTService([]byte("test"), 30, "localhost", true)
	s := newServer(store, jwtService)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    u.GetEmail(),
				"password": u.GetPassword(),
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email": "invalid",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			b := &bytes.Buffer{}
			assert.NoError(t, json.NewEncoder(b).Encode(tc.payload))

			req, _ := http.NewRequest(http.MethodPost, "/login", b)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedCode == http.StatusOK {
				cookies := rec.Result().Cookies()
				assert.NotEmpty(t, cookies)

				var jwtCookie *http.Cookie
				for _, c := range cookies {
					if c.Name == "jwt-token" {
						jwtCookie = c
						break
					}
				}
				assert.NotNil(t, jwtCookie)
				assert.Equal(t, jwtCookie.Name, "jwt-token")
				assert.NotEmpty(t, jwtCookie.Value)
			}
		})
	}
}

func TestServer_HandleUserLogout(t *testing.T) {
	u := domain.TestUser(t, "email@example.com", "a password")

	store := teststore.New()
	assert.NoError(t, store.User().SaveUser(context.Background(), u))

	jwtService := jwt.NewJWTService([]byte("test"), 30, "localhost", true)
	s := newServer(store, jwtService)

	rec := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodPost, "/logout", nil)

	s.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)

	cookies := rec.Result().Cookies()
	assert.NotEmpty(t, cookies)

	var jwtCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "jwt-token" {
			jwtCookie = c
			break
		}
	}
	assert.NotNil(t, jwtCookie)
	assert.Equal(t, jwtCookie.Name, "jwt-token")
	assert.Empty(t, jwtCookie.Value)
	assert.Equal(t, jwtCookie.MaxAge, -1)
}

func TestServer_JWTProtectedMiddleware(t *testing.T) {
	u := domain.TestUser(t, "email@example.com", "a password")

	store := teststore.New()
	assert.NoError(t, store.User().SaveUser(context.Background(), u))

	testCases := []struct {
		name         string
		userID       uuid.UUID
		expectedCode int
	}{
		{
			name:         "authenticated",
			userID:       u.GetID(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "not authenticated",
			userID:       uuid.Nil,
			expectedCode: http.StatusUnauthorized,
		},
	}

	jwtService := jwt.NewJWTService([]byte("test"), 30, "localhost", true)
	s := newServer(store, jwtService)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)

			token, _ := jwtService.CreateJWTTokenForUser(tc.userID)

			req.Header.Set("Cookie", jwtService.AuthCookie(token).String())

			s.jwtProtectedMiddleware(handler).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleTasksCreate(t *testing.T) {
	store := teststore.New()
	jwtService := jwt.NewJWTService([]byte("test"), 30, "localhost", true)
	s := newServer(store, jwtService)

	u := domain.TestUser(t, "email@example.com", "a password")
	assert.NoError(t, store.User().SaveUser(context.Background(), u))
	token, _ := jwtService.CreateJWTTokenForUser(u.GetID())

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"task_text":  "test task",
				"task_order": 1,
				"user_id":    u.GetID(),
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]interface{}{
				"task_text": "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			b := &bytes.Buffer{}
			assert.NoError(t, json.NewEncoder(b).Encode(tc.payload))

			url := fmt.Sprintf("/users/%s/tasks", u.GetID())
			req, _ := http.NewRequest(http.MethodPost, url, b)
			req.Header.Set("Cookie", jwtService.AuthCookie(token).String())

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleTasksGetAllByUser(t *testing.T) {
	store := teststore.New()
	jwtService := jwt.NewJWTService([]byte("test"), 30, "localhost", true)
	s := newServer(store, jwtService)

	u := domain.TestUser(t, "email@example.com", "a password")
	assert.NoError(t, store.User().SaveUser(context.Background(), u))

	task := domain.TestTask(t, "test task", 0, false, u.GetID())
	assert.NoError(t, store.Task().SaveTask(context.Background(), task))

	token, _ := jwtService.CreateJWTTokenForUser(u.GetID())

	rec := httptest.NewRecorder()

	url := fmt.Sprintf("/users/%s/tasks", u.GetID())
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Cookie", jwtService.AuthCookie(token).String())

	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestServer_HandleTasksMarkAsDone(t *testing.T) {
	store := teststore.New()
	jwtService := jwt.NewJWTService([]byte("test"), 30, "localhost", true)
	s := newServer(store, jwtService)

	u := domain.TestUser(t, "email@example.com", "a password")
	assert.NoError(t, store.User().SaveUser(context.Background(), u))

	task := domain.TestTask(t, "test task", 0, false, u.GetID())
	assert.NoError(t, store.Task().SaveTask(context.Background(), task))

	token, _ := jwtService.CreateJWTTokenForUser(u.GetID())

	rec := httptest.NewRecorder()

	url := fmt.Sprintf("/users/%s/tasks/%s/mark-done", u.GetID(), task.GetID())
	req, _ := http.NewRequest(http.MethodPut, url, nil)
	req.Header.Set("Cookie", jwtService.AuthCookie(token).String())

	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestServer_HandleTasksMarkAsNotDone(t *testing.T) {
	store := teststore.New()
	jwtService := jwt.NewJWTService([]byte("test"), 30, "localhost", true)
	s := newServer(store, jwtService)

	u := domain.TestUser(t, "email@example.com", "a password")
	assert.NoError(t, store.User().SaveUser(context.Background(), u))

	task := domain.TestTask(t, "test task", 0, true, u.GetID())
	assert.NoError(t, store.Task().SaveTask(context.Background(), task))

	token, _ := jwtService.CreateJWTTokenForUser(u.GetID())

	rec := httptest.NewRecorder()

	b := &bytes.Buffer{}
	err := json.NewEncoder(b).Encode(map[string]interface{}{
		"task_text":  "updated task text",
		"task_order": 10,
	})
	assert.NoError(t, err)

	url := fmt.Sprintf("/users/%s/tasks/%s", u.GetID(), task.GetID())
	req, _ := http.NewRequest(http.MethodPut, url, b)
	req.Header.Set("Cookie", jwtService.AuthCookie(token).String())

	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestServer_HandleTasksDelete(t *testing.T) {
	store := teststore.New()
	jwtService := jwt.NewJWTService([]byte("test"), 30, "localhost", true)
	s := newServer(store, jwtService)

	u := domain.TestUser(t, "email@example.com", "a password")
	assert.NoError(t, store.User().SaveUser(context.Background(), u))

	task := domain.TestTask(t, "test task", 0, true, u.GetID())
	assert.NoError(t, store.Task().SaveTask(context.Background(), task))

	token, _ := jwtService.CreateJWTTokenForUser(u.GetID())

	rec := httptest.NewRecorder()

	url := fmt.Sprintf("/users/%s/tasks/%s", u.GetID(), task.GetID())
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Cookie", jwtService.AuthCookie(token).String())

	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}
