package jwt_test

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/jwt"
)

var testKey = []byte("0123456789abcdef0123456789abcdef")

func newService(t *testing.T, mutate ...func(*jwt.Config)) *jwt.Service {
	t.Helper()
	cfg := jwt.Config{SigningKey: testKey, Issuer: "todo-api", TTL: 30 * time.Minute, SecureCookie: true}
	for _, m := range mutate {
		m(&cfg)
	}
	s, err := jwt.NewService(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func sign(t *testing.T, method gojwt.SigningMethod, key any, claims gojwt.Claims) string {
	t.Helper()
	s, err := gojwt.NewWithClaims(method, claims).SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func validClaims(sub string) gojwt.RegisteredClaims {
	now := time.Now()
	return gojwt.RegisteredClaims{
		Subject: sub, Issuer: "todo-api",
		IssuedAt: gojwt.NewNumericDate(now), ExpiresAt: gojwt.NewNumericDate(now.Add(time.Hour)),
	}
}

func TestNewService_RejectsWeakConfig(t *testing.T) {
	t.Parallel()
	for name, cfg := range map[string]jwt.Config{
		"short key": {SigningKey: []byte("go-todo-app"), Issuer: "x", TTL: time.Minute},
		"no ttl":    {SigningKey: testKey, Issuer: "x"},
		"no issuer": {SigningKey: testKey, TTL: time.Minute},
	} {
		if _, err := jwt.NewService(cfg); err == nil {
			t.Errorf("%s: NewService succeeded, want error", name)
		}
	}
}

func TestService_IssueAndParse_RoundTrip(t *testing.T) {
	t.Parallel()
	s := newService(t)
	id := uuid.New()

	token, expires, err := s.Issue(id)
	if err != nil {
		t.Fatal(err)
	}
	if d := time.Until(expires); d < 29*time.Minute || d > 31*time.Minute {
		t.Fatalf("expires in %v, want ~30m", d)
	}

	got, err := s.Parse(token)
	if err != nil || got != id {
		t.Fatalf("Parse = %v, %v; want %v", got, err, id)
	}
}

func TestService_Parse_RejectsBadTokens(t *testing.T) {
	t.Parallel()
	s := newService(t)
	sub := uuid.NewString()
	valid, _, _ := s.Issue(uuid.MustParse(sub))

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"` + sub + `","iss":"todo-api","exp":9999999999}`))

	expired := validClaims(sub)
	expired.ExpiresAt = gojwt.NewNumericDate(time.Now().Add(-time.Hour))
	noExp := validClaims(sub)
	noExp.ExpiresAt = nil
	wrongIss := validClaims(sub)
	wrongIss.Issuer = "someone-else"
	future := validClaims(sub)
	future.IssuedAt = gojwt.NewNumericDate(time.Now().Add(time.Hour))

	tests := map[string]string{
		"empty":            "",
		"garbage":          "garbage",
		"two segments":     "a.b",
		"tampered":         valid[:len(valid)-2] + "xx",
		"alg none":         header + "." + payload + ".",
		"HS512 same key":   sign(t, gojwt.SigningMethodHS512, testKey, validClaims(sub)),
		"wrong key":        sign(t, gojwt.SigningMethodHS256, []byte(strings.Repeat("k", 32)), validClaims(sub)),
		"expired":          sign(t, gojwt.SigningMethodHS256, testKey, expired),
		"missing exp":      sign(t, gojwt.SigningMethodHS256, testKey, noExp),
		"wrong issuer":     sign(t, gojwt.SigningMethodHS256, testKey, wrongIss),
		"issued in future": sign(t, gojwt.SigningMethodHS256, testKey, future),
		"non-uuid subject": sign(t, gojwt.SigningMethodHS256, testKey, validClaims("alice")),
		"nil-uuid subject": sign(t, gojwt.SigningMethodHS256, testKey, validClaims(uuid.Nil.String())),
	}

	for name, token := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if _, err := s.Parse(token); !errors.Is(err, jwt.ErrInvalidToken) {
				t.Fatalf("err = %v, want ErrInvalidToken", err)
			}
		})
	}
}

func TestService_UserIDFromRequest(t *testing.T) {
	t.Parallel()
	s := newService(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if _, err := s.UserIDFromRequest(req); !errors.Is(err, jwt.ErrNoToken) {
		t.Fatalf("no cookie: err = %v, want ErrNoToken", err)
	}

	id := uuid.New()
	token, _, _ := s.Issue(id)
	req.AddCookie(&http.Cookie{Name: jwt.CookieName, Value: token})
	if got, err := s.UserIDFromRequest(req); err != nil || got != id {
		t.Fatalf("got %v, %v", got, err)
	}
}

func TestService_Cookies_HonourConfig(t *testing.T) {
	t.Parallel()

	for _, secure := range []bool{true, false} {
		s := newService(t, func(c *jwt.Config) { c.SecureCookie = secure; c.CookieDomain = "example.com" })

		w := httptest.NewRecorder()
		s.SetCookie(w, "tok", time.Now().Add(30*time.Minute))
		c := w.Result().Cookies()[0]
		if c.Name != jwt.CookieName || c.Value != "tok" || !c.HttpOnly || c.Secure != secure ||
			c.SameSite != http.SameSiteLaxMode || c.MaxAge != 1800 || c.Domain != "example.com" || c.Path != "/" {
			t.Fatalf("secure=%v: unexpected cookie %+v", secure, c)
		}

		w = httptest.NewRecorder()
		s.ClearCookie(w)
		c = w.Result().Cookies()[0]
		if c.Value != "" || c.MaxAge >= 0 || c.Secure != secure || !c.HttpOnly {
			t.Fatalf("secure=%v: unexpected cleared cookie %+v", secure, c)
		}
	}
}
