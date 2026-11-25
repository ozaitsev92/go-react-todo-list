package v1_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	v1 "github.com/ozaitsev92/tododdd/internal/controller/http/v1"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/jwt"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/model"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/password"
	taskmemory "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/memory"
	usermemory "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/memory"
	"github.com/ozaitsev92/tododdd/internal/usecase"
	"github.com/ozaitsev92/tododdd/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

const (
	testOrigin   = "http://localhost:8081"
	testPassword = "Secret123"
)

func init() { gin.SetMode(gin.TestMode) }

type env struct {
	t      *testing.T
	router http.Handler
	tokens *jwt.Service
	tasks  task.Repository
}

type option func(*v1.Deps, *env)

func withTaskRepo(r task.Repository) option {
	return func(_ *v1.Deps, e *env) { e.tasks = r }
}

func withReady(f func(context.Context) error) option {
	return func(d *v1.Deps, _ *env) { d.Ready = f }
}

func newEnv(t *testing.T, opts ...option) *env {
	t.Helper()

	e := &env{t: t, tasks: taskmemory.NewRepository()}
	hasher, err := password.NewBcryptHasher(bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	users, err := usecase.NewUserUseCase(usermemory.NewRepository(), hasher)
	if err != nil {
		t.Fatal(err)
	}
	tokens, err := jwt.NewService(jwt.Config{
		SigningKey: []byte(strings.Repeat("k", 32)), Issuer: "todo-api", TTL: 30 * time.Minute,
	})
	if err != nil {
		t.Fatal(err)
	}

	d := v1.Deps{
		Log: logger.Discard(), Tokens: tokens, Users: users,
		AllowedOrigin: testOrigin, MaxBodyBytes: 4 << 10,
	}
	for _, o := range opts {
		o(&d, e)
	}
	d.Tasks = usecase.NewTaskUseCase(e.tasks)

	e.router = v1.NewRouter(d)
	e.tokens = tokens

	return e
}

type response struct {
	*httptest.ResponseRecorder
}

func (r response) errorMessage(t *testing.T) string {
	t.Helper()
	var body struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(r.Body.Bytes(), &body); err != nil || body.Error == "" {
		t.Fatalf("body %q is not an error response: %v", r.Body, err)
	}
	return body.Error
}

func (e *env) do(method, path, body string, cookie *http.Cookie) response {
	e.t.Helper()
	var rd io.Reader
	if body != "" {
		rd = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, rd)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	e.router.ServeHTTP(w, req)
	return response{w}
}

// signup registers and logs in a user, returning their session cookie and ID.
func (e *env) signup(email string) (*http.Cookie, string) {
	e.t.Helper()
	creds := `{"email":"` + email + `","password":"` + testPassword + `"}`
	if r := e.do(http.MethodPost, "/v1/users", creds, nil); r.Code != http.StatusCreated {
		e.t.Fatalf("signup %s: %d %s", email, r.Code, r.Body)
	}
	r := e.do(http.MethodPost, "/v1/users/login", creds, nil)
	if r.Code != http.StatusOK {
		e.t.Fatalf("login %s: %d %s", email, r.Code, r.Body)
	}
	var u model.User
	decode(e.t, r, &u)
	for _, c := range r.Result().Cookies() {
		if c.Name == jwt.CookieName {
			return &http.Cookie{Name: c.Name, Value: c.Value}, u.ID
		}
	}
	e.t.Fatal("login set no session cookie")
	return nil, ""
}

func (e *env) createTask(cookie *http.Cookie, text string) model.Task {
	e.t.Helper()
	r := e.do(http.MethodPost, "/v1/tasks", `{"text":"`+text+`"}`, cookie)
	if r.Code != http.StatusCreated {
		e.t.Fatalf("create task: %d %s", r.Code, r.Body)
	}
	var tk model.Task
	decode(e.t, r, &tk)
	return tk
}

func decode(t *testing.T, r response, v any) {
	t.Helper()
	if err := json.Unmarshal(r.Body.Bytes(), v); err != nil {
		t.Fatalf("decode %q: %v", r.Body, err)
	}
}

func TestOps_HealthReadyMetrics(t *testing.T) {
	t.Parallel()

	e := newEnv(t)
	if r := e.do(http.MethodGet, "/healthz", "", nil); r.Code != http.StatusOK {
		t.Fatalf("/healthz = %d", r.Code)
	}
	if r := e.do(http.MethodGet, "/readyz", "", nil); r.Code != http.StatusOK {
		t.Fatalf("/readyz = %d", r.Code)
	}
	e.do(http.MethodGet, "/v1/tasks", "", nil)
	r := e.do(http.MethodGet, "/metrics", "", nil)
	if r.Code != http.StatusOK || !strings.Contains(r.Body.String(), `http_requests_total{method="GET",route="/v1/tasks",status="401"} 1`) {
		t.Fatalf("/metrics = %d, missing request metric:\n%s", r.Code, r.Body)
	}

	down := newEnv(t, withReady(func(context.Context) error { return errors.New("mongo down") }))
	if r := down.do(http.MethodGet, "/readyz", "", nil); r.Code != http.StatusServiceUnavailable {
		t.Fatalf("/readyz with db down = %d, want 503", r.Code)
	}
}

func TestRouter_UnknownRouteAndMethod(t *testing.T) {
	t.Parallel()
	e := newEnv(t)
	if r := e.do(http.MethodGet, "/nope", "", nil); r.Code != http.StatusNotFound {
		t.Fatalf("unknown route = %d", r.Code)
	}
	if r := e.do(http.MethodPatch, "/v1/users", "", nil); r.Code != http.StatusMethodNotAllowed {
		t.Fatalf("wrong method = %d", r.Code)
	}
}

func TestUserHandler_Register(t *testing.T) {
	t.Parallel()

	t.Run("success never returns password", func(t *testing.T) {
		t.Parallel()
		e := newEnv(t)
		r := e.do(http.MethodPost, "/v1/users", `{"email":"Alice@Example.com","password":"`+testPassword+`"}`, nil)
		if r.Code != http.StatusCreated {
			t.Fatalf("status %d %s", r.Code, r.Body)
		}
		if strings.Contains(strings.ToLower(r.Body.String()), "password") {
			t.Fatalf("response leaks password: %s", r.Body)
		}
		var u model.User
		decode(t, r, &u)
		if u.Email != "alice@example.com" || uuid.Validate(u.ID) != nil {
			t.Fatalf("got %+v", u)
		}
	})

	tests := []struct {
		name, body string
		want       int
	}{
		{"malformed JSON", `{"email":`, http.StatusBadRequest},
		{"missing password", `{"email":"a@example.com"}`, http.StatusBadRequest},
		{"operator injection", `{"email":{"$ne":""},"password":"x"}`, http.StatusBadRequest},
		{"invalid email", `{"email":"nope","password":"` + testPassword + `"}`, http.StatusUnprocessableEntity},
		{"weak password", `{"email":"a@example.com","password":"a"}`, http.StatusUnprocessableEntity},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newEnv(t).do(http.MethodPost, "/v1/users", tt.body, nil)
			if r.Code != tt.want {
				t.Fatalf("status %d, want %d (%s)", r.Code, tt.want, r.Body)
			}
			r.errorMessage(t)
		})
	}

	t.Run("duplicate in any case returns 409", func(t *testing.T) {
		t.Parallel()
		e := newEnv(t)
		e.signup("bob@example.com")
		r := e.do(http.MethodPost, "/v1/users", `{"email":"BOB@example.com","password":"`+testPassword+`"}`, nil)
		if r.Code != http.StatusConflict {
			t.Fatalf("status %d, want 409", r.Code)
		}
	})
}

func TestUserHandler_LoginLogoutCurrent(t *testing.T) {
	t.Parallel()
	e := newEnv(t)
	cookie, id := e.signup("alice@example.com")

	t.Run("bad credentials are indistinguishable", func(t *testing.T) {
		t.Parallel()
		wrong := e.do(http.MethodPost, "/v1/users/login", `{"email":"alice@example.com","password":"wrong-password"}`, nil)
		unknown := e.do(http.MethodPost, "/v1/users/login", `{"email":"ghost@example.com","password":"wrong-password"}`, nil)
		if wrong.Code != http.StatusUnauthorized || unknown.Code != http.StatusUnauthorized ||
			wrong.Body.String() != unknown.Body.String() {
			t.Fatalf("wrong=%d %s unknown=%d %s", wrong.Code, wrong.Body, unknown.Code, unknown.Body)
		}
		if len(wrong.Result().Cookies()) != 0 {
			t.Fatal("failed login set a cookie")
		}
	})

	t.Run("current returns the session user", func(t *testing.T) {
		t.Parallel()
		r := e.do(http.MethodGet, "/v1/users/current", "", cookie)
		var u model.User
		decode(t, r, &u)
		if r.Code != http.StatusOK || u.ID != id {
			t.Fatalf("status %d user %+v", r.Code, u)
		}
	})

	t.Run("current without session is 401", func(t *testing.T) {
		t.Parallel()
		if r := e.do(http.MethodGet, "/v1/users/current", "", nil); r.Code != http.StatusUnauthorized {
			t.Fatalf("status %d", r.Code)
		}
	})

	t.Run("logout clears the cookie", func(t *testing.T) {
		t.Parallel()
		r := e.do(http.MethodPost, "/v1/users/logout", "", cookie)
		if r.Code != http.StatusNoContent {
			t.Fatalf("status %d", r.Code)
		}
		cs := r.Result().Cookies()
		if len(cs) != 1 || cs[0].Name != jwt.CookieName || cs[0].MaxAge >= 0 {
			t.Fatalf("cookies %+v", cs)
		}
	})
}

func TestTaskHandler_RequiresValidSession(t *testing.T) {
	t.Parallel()
	e := newEnv(t)
	ghost, _, _ := e.tokens.Issue(uuid.New()) // valid signature, user does not exist

	for name, c := range map[string]*http.Cookie{
		"no cookie":    nil,
		"garbage":      {Name: jwt.CookieName, Value: "garbage"},
		"unknown user": {Name: jwt.CookieName, Value: ghost},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r := e.do(http.MethodGet, "/v1/tasks", "", c)
			if r.Code != http.StatusUnauthorized {
				t.Fatalf("status %d, want 401 (%s)", r.Code, r.Body)
			}
			r.errorMessage(t)
		})
	}
}

func TestTaskHandler_CRUD(t *testing.T) {
	t.Parallel()
	e := newEnv(t)
	alice, aliceID := e.signup("alice@example.com")

	list := func() []model.Task {
		r := e.do(http.MethodGet, "/v1/tasks", "", alice)
		if r.Code != http.StatusOK {
			t.Fatalf("list: %d", r.Code)
		}
		var ts []model.Task
		decode(t, r, &ts)
		return ts
	}

	if ts := list(); ts == nil || len(ts) != 0 {
		t.Fatalf("empty list = %#v, want []", ts)
	}

	r := e.do(http.MethodPost, "/v1/tasks", `{"text":"  buy milk "}`, alice)
	if r.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", r.Code, r.Body)
	}
	var created model.Task
	decode(t, r, &created)
	if created.Text != "buy milk" || created.UserID != aliceID || created.Completed {
		t.Fatalf("created %+v", created)
	}
	if loc := r.Header().Get("Location"); loc != "/v1/tasks/"+created.ID {
		t.Fatalf("Location = %q", loc)
	}

	r = e.do(http.MethodGet, "/v1/tasks/"+created.ID, "", alice)
	if r.Code != http.StatusOK {
		t.Fatalf("get: %d", r.Code)
	}

	r = e.do(http.MethodPut, "/v1/tasks/"+created.ID, `{"text":"buy oat milk"}`, alice)
	var updated model.Task
	decode(t, r, &updated)
	if r.Code != http.StatusOK || updated.Text != "buy oat milk" {
		t.Fatalf("update: %d %+v", r.Code, updated)
	}

	for path, want := range map[string]bool{"/mark-completed": true, "/mark-not-completed": false} {
		r = e.do(http.MethodPut, "/v1/tasks/"+created.ID+path, "", alice)
		var tk model.Task
		decode(t, r, &tk)
		if r.Code != http.StatusOK || tk.Completed != want {
			t.Fatalf("%s: %d %+v", path, r.Code, tk)
		}
	}

	second := e.createTask(alice, "second")
	if ts := list(); len(ts) != 2 || ts[0].ID != created.ID || ts[1].ID != second.ID {
		t.Fatalf("list order %+v", ts)
	}

	if r = e.do(http.MethodDelete, "/v1/tasks/"+created.ID, "", alice); r.Code != http.StatusNoContent || r.Body.Len() != 0 {
		t.Fatalf("delete: %d %q", r.Code, r.Body)
	}
	if r = e.do(http.MethodDelete, "/v1/tasks/"+created.ID, "", alice); r.Code != http.StatusNotFound {
		t.Fatalf("delete twice: %d", r.Code)
	}
	if ts := list(); len(ts) != 1 {
		t.Fatalf("after delete %+v", ts)
	}
}

func TestTaskHandler_InputValidation(t *testing.T) {
	t.Parallel()
	e := newEnv(t)
	alice, _ := e.signup("alice@example.com")
	tk := e.createTask(alice, "x")

	tests := []struct {
		name, method, path, body string
		want                     int
	}{
		{"create malformed JSON", http.MethodPost, "/v1/tasks", `{"text":`, http.StatusBadRequest},
		{"create whitespace text", http.MethodPost, "/v1/tasks", `{"text":"   "}`, http.StatusUnprocessableEntity},
		{"create oversized body", http.MethodPost, "/v1/tasks", `{"text":"` + strings.Repeat("a", 8<<10) + `"}`, http.StatusRequestEntityTooLarge},
		{"update invalid id", http.MethodPut, "/v1/tasks/not-a-uuid", `{"text":"y"}`, http.StatusBadRequest},
		{"update invalid text", http.MethodPut, "/v1/tasks/" + tk.ID, `{"text":" "}`, http.StatusUnprocessableEntity},
		{"get invalid id", http.MethodGet, "/v1/tasks/not-a-uuid", "", http.StatusBadRequest},
		{"delete invalid id", http.MethodDelete, "/v1/tasks/not-a-uuid", "", http.StatusBadRequest},
		{"mark invalid id", http.MethodPut, "/v1/tasks/not-a-uuid/mark-completed", "", http.StatusBadRequest},
		{"update missing", http.MethodPut, "/v1/tasks/" + uuid.NewString(), `{"text":"y"}`, http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := e.do(tt.method, tt.path, tt.body, alice)
			if r.Code != tt.want {
				t.Fatalf("status %d, want %d (%s)", r.Code, tt.want, r.Body)
			}
			r.errorMessage(t)
		})
	}
}

func TestTaskHandler_OtherUsersTask_Returns404AndIsUntouched(t *testing.T) {
	t.Parallel()
	e := newEnv(t)
	alice, _ := e.signup("alice@example.com")
	bob, _ := e.signup("bob@example.com")
	tk := e.createTask(alice, "alice's secret")

	attempts := []struct{ method, path, body string }{
		{http.MethodGet, "/v1/tasks/" + tk.ID, ""},
		{http.MethodPut, "/v1/tasks/" + tk.ID, `{"text":"pwned"}`},
		{http.MethodPut, "/v1/tasks/" + tk.ID + "/mark-completed", ""},
		{http.MethodPut, "/v1/tasks/" + tk.ID + "/mark-not-completed", ""},
		{http.MethodDelete, "/v1/tasks/" + tk.ID, ""},
	}
	missing := e.do(http.MethodGet, "/v1/tasks/"+uuid.NewString(), "", bob)

	for _, a := range attempts {
		r := e.do(a.method, a.path, a.body, bob)
		if r.Code != http.StatusNotFound || r.Body.String() != missing.Body.String() {
			t.Errorf("%s %s by bob: %d %s (must look exactly like a missing task)", a.method, a.path, r.Code, r.Body)
		}
	}

	var bobs []model.Task
	decode(t, e.do(http.MethodGet, "/v1/tasks", "", bob), &bobs)
	if len(bobs) != 0 {
		t.Fatalf("bob sees %+v", bobs)
	}

	r := e.do(http.MethodGet, "/v1/tasks/"+tk.ID, "", alice)
	var got model.Task
	decode(t, r, &got)
	if got.Text != "alice's secret" || got.Completed {
		t.Fatalf("alice's task was modified: %+v", got)
	}
}

// failingTasks fails GetAllByUserID and Save, to check 500 handling.
type failingTasks struct{ task.Repository }

var errInternal = errors.New("mongo: connection refused at 10.0.0.5:27017")

func (failingTasks) GetAllByUserID(context.Context, uuid.UUID) ([]task.Task, error) {
	return nil, errInternal
}
func (failingTasks) Save(context.Context, task.Task) error { return errInternal }

func TestTaskHandler_InternalError_DoesNotLeakDetails(t *testing.T) {
	t.Parallel()
	e := newEnv(t, withTaskRepo(failingTasks{taskmemory.NewRepository()}))
	alice, _ := e.signup("alice@example.com")

	for _, r := range []response{
		e.do(http.MethodGet, "/v1/tasks", "", alice),
		e.do(http.MethodPost, "/v1/tasks", `{"text":"x"}`, alice),
	} {
		if r.Code != http.StatusInternalServerError {
			t.Fatalf("status %d", r.Code)
		}
		if msg := r.errorMessage(t); msg != "internal server error" {
			t.Fatalf("leaked error message %q", msg)
		}
	}
}

func TestCORS_ThroughRouter(t *testing.T) {
	t.Parallel()
	e := newEnv(t)

	req := httptest.NewRequest(http.MethodOptions, "/v1/tasks/"+uuid.NewString(), nil)
	req.Header.Set("Origin", testOrigin)
	req.Header.Set("Access-Control-Request-Method", http.MethodDelete)
	w := httptest.NewRecorder()
	e.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent || w.Header().Get("Access-Control-Allow-Origin") != testOrigin {
		t.Fatalf("preflight %d %v", w.Code, w.Header())
	}
}
