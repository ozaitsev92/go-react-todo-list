package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/middleware"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	"github.com/ozaitsev92/tododdd/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func init() { gin.SetMode(gin.TestMode) }

type stubTokens struct {
	id  uuid.UUID
	err error
}

func (s stubTokens) UserIDFromRequest(*http.Request) (uuid.UUID, error) { return s.id, s.err }

type stubUsers struct {
	u   user.User
	err error
}

func (s stubUsers) GetUserByID(context.Context, uuid.UUID) (user.User, error) { return s.u, s.err }

func serve(h http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func TestAuthenticate(t *testing.T) {
	t.Parallel()
	alice := user.User{ID: uuid.New(), Email: "alice@example.com"}

	tests := []struct {
		name       string
		tokens     stubTokens
		users      stubUsers
		wantStatus int
	}{
		{"valid", stubTokens{id: alice.ID}, stubUsers{u: alice}, http.StatusOK},
		{"invalid token", stubTokens{err: errors.New("bad")}, stubUsers{u: alice}, http.StatusUnauthorized},
		{"deleted user", stubTokens{id: alice.ID}, stubUsers{err: user.ErrUserNotFound}, http.StatusUnauthorized},
		{"store down", stubTokens{id: alice.ID}, stubUsers{err: errors.New("db down")}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := gin.New()
			r.GET("/", middleware.Authenticate(tt.tokens, tt.users, logger.Discard()), func(c *gin.Context) {
				u, ok := middleware.CurrentUser(c.Request.Context())
				if !ok || u.ID != alice.ID {
					t.Errorf("CurrentUser = %+v, %v", u, ok)
				}
				c.Status(http.StatusOK)
			})

			w := serve(r, httptest.NewRequest(http.MethodGet, "/", nil))
			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body %s)", w.Code, tt.wantStatus, w.Body)
			}
			if w.Code >= 400 && !strings.Contains(w.Body.String(), `"error"`) {
				t.Fatalf("body %s lacks error field", w.Body)
			}
		})
	}
}

func TestBodyLimit(t *testing.T) {
	t.Parallel()
	r := gin.New()
	r.Use(middleware.BodyLimit(10))
	r.POST("/", func(c *gin.Context) {
		var v map[string]any
		if err := c.ShouldBindJSON(&v); err != nil {
			if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
				c.Status(http.StatusRequestEntityTooLarge)
				return
			}
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusOK)
	})

	if w := serve(r, httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"a":1}`))); w.Code != http.StatusOK {
		t.Fatalf("small body status %d", w.Code)
	}
	big := `{"a":"` + strings.Repeat("x", 100) + `"}`
	if w := serve(r, httptest.NewRequest(http.MethodPost, "/", strings.NewReader(big))); w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("big body status %d", w.Code)
	}
}

func TestMetrics_UsesRouteTemplate(t *testing.T) {
	t.Parallel()
	reg := prometheus.NewRegistry()
	m := middleware.NewMetrics(reg)

	r := gin.New()
	r.Use(m.Handler())
	r.GET("/tasks/:id", func(c *gin.Context) { c.Status(http.StatusOK) })

	for range 3 {
		serve(r, httptest.NewRequest(http.MethodGet, "/tasks/"+uuid.NewString(), nil))
	}
	serve(r, httptest.NewRequest(http.MethodGet, "/nope", nil))

	expected := `
# HELP http_requests_total HTTP requests by method, route and status.
# TYPE http_requests_total counter
http_requests_total{method="GET",route="/tasks/:id",status="200"} 3
http_requests_total{method="GET",route="unmatched",status="404"} 1
`
	if err := testutil.GatherAndCompare(reg, strings.NewReader(expected), "http_requests_total"); err != nil {
		t.Fatal(err)
	}
}

func TestRecovery_ReturnsJSON500(t *testing.T) {
	t.Parallel()
	r := gin.New()
	r.Use(middleware.Recovery(logger.Discard()))
	r.GET("/", func(*gin.Context) { panic("boom") })

	w := serve(r, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusInternalServerError || !strings.Contains(w.Body.String(), "internal server error") ||
		strings.Contains(w.Body.String(), "boom") {
		t.Fatalf("status %d body %s", w.Code, w.Body)
	}
}

func TestMiddleware_PanicIsLoggedAndCounted(t *testing.T) {
	t.Parallel()
	reg := prometheus.NewRegistry()
	var logs strings.Builder

	r := gin.New()
	r.Use(
		middleware.RequestLogger(logger.New("info", &logs)),
		middleware.NewMetrics(reg).Handler(),
		middleware.Recovery(logger.Discard()),
	)
	r.GET("/boom", func(*gin.Context) { panic("boom") })

	if w := serve(r, httptest.NewRequest(http.MethodGet, "/boom", nil)); w.Code != http.StatusInternalServerError {
		t.Fatalf("status %d, want 500", w.Code)
	}

	expected := `
# HELP http_requests_total HTTP requests by method, route and status.
# TYPE http_requests_total counter
http_requests_total{method="GET",route="/boom",status="500"} 1
`
	if err := testutil.GatherAndCompare(reg, strings.NewReader(expected), "http_requests_total"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(logs.String(), `"status":500`) {
		t.Fatalf("panicked request not logged: %s", logs.String())
	}
}

func TestRequestLogger_PropagatesOrGeneratesRequestID(t *testing.T) {
	t.Parallel()
	r := gin.New()
	r.Use(middleware.RequestLogger(logger.Discard()))
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(middleware.RequestIDHeader, "abc-123")
	if got := serve(r, req).Header().Get(middleware.RequestIDHeader); got != "abc-123" {
		t.Fatalf("request id = %q, want abc-123", got)
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(middleware.RequestIDHeader, "bad id\nwith newline")
	if got := serve(r, req).Header().Get(middleware.RequestIDHeader); uuid.Validate(got) != nil {
		t.Fatalf("request id = %q, want generated uuid", got)
	}
}
