// Package jwt issues and verifies session tokens and manages their cookie.
// Tokens are stateless HS256 JWTs and cannot be revoked before they expire.
package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CookieName is the name of the session cookie.
const CookieName = "jwt-token"

// MinKeyLength is the minimum HMAC key length in bytes.
const MinKeyLength = 32

var (
	// ErrNoToken is returned when the request carries no session cookie.
	ErrNoToken = errors.New("jwt: no token")
	// ErrInvalidToken is returned for malformed, expired or forged tokens.
	ErrInvalidToken = errors.New("jwt: invalid token")
)

// Config configures a Service.
type Config struct {
	SigningKey   []byte
	Issuer       string
	TTL          time.Duration
	CookieDomain string
	SecureCookie bool
}

// Service issues and verifies tokens.
type Service struct {
	cfg    Config
	parser *jwt.Parser
	now    func() time.Time
}

// NewService validates cfg and returns a Service.
func NewService(cfg Config) (*Service, error) {
	if len(cfg.SigningKey) < MinKeyLength {
		return nil, fmt.Errorf("jwt: signing key must be at least %d bytes", MinKeyLength)
	}
	if cfg.TTL <= 0 {
		return nil, errors.New("jwt: TTL must be positive")
	}
	if cfg.Issuer == "" {
		return nil, errors.New("jwt: issuer is required")
	}

	return &Service{
		cfg: cfg,
		parser: jwt.NewParser(
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
			jwt.WithIssuer(cfg.Issuer),
			jwt.WithExpirationRequired(),
			jwt.WithIssuedAt(),
			jwt.WithLeeway(5*time.Second),
		),
		now: time.Now,
	}, nil
}

// Issue creates a signed token for userID and returns it with its expiry.
func (s *Service) Issue(userID uuid.UUID) (string, time.Time, error) {
	now := s.now()
	expires := now.Add(s.cfg.TTL)

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		Issuer:    s.cfg.Issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expires),
		ID:        uuid.NewString(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.cfg.SigningKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("jwt: sign: %w", err)
	}

	return token, expires, nil
}

// Parse verifies token and returns the user ID it was issued for.
func (s *Service) Parse(token string) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims

	parsed, err := s.parser.ParseWithClaims(token, &claims, func(*jwt.Token) (any, error) {
		return s.cfg.SigningKey, nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}
	if !parsed.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil || userID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%w: bad subject", ErrInvalidToken)
	}

	return userID, nil
}

// UserIDFromRequest extracts and verifies the session cookie.
func (s *Service) UserIDFromRequest(r *http.Request) (uuid.UUID, error) {
	c, err := r.Cookie(CookieName)
	if err != nil || c.Value == "" {
		return uuid.Nil, ErrNoToken
	}

	return s.Parse(c.Value)
}

// SetCookie writes the session cookie for token.
func (s *Service) SetCookie(w http.ResponseWriter, token string, expires time.Time) {
	c := s.cookie() //nolint:gosec // see cookie
	c.Value = token
	c.Expires = expires
	c.MaxAge = int(s.cfg.TTL.Seconds())
	http.SetCookie(w, c)
}

// ClearCookie expires the session cookie.
func (s *Service) ClearCookie(w http.ResponseWriter) {
	c := s.cookie() //nolint:gosec // see cookie
	c.MaxAge = -1
	http.SetCookie(w, c)
}

// cookie returns the base cookie. Secure is configurable for local plain HTTP.
func (s *Service) cookie() *http.Cookie {
	return &http.Cookie{ //nolint:gosec // Secure comes from config
		Name:     CookieName,
		Path:     "/",
		Domain:   s.cfg.CookieDomain,
		Secure:   s.cfg.SecureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}
