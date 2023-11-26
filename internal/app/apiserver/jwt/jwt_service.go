package jwt

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	jwtCookieName = "jwt-token"
)

var (
	errSigningMethodMismatch = errors.New("signing method mismatch")
)

type JWTService struct {
	jwtSigningKey    []byte
	defaultCookie    http.Cookie
	jwtSessionLength time.Duration
	jwtSigningMethod *jwt.SigningMethodHMAC
}

func NewJWTService(signingKey []byte, sessionLength int, cookieDomain string, secureCookie bool) *JWTService {
	return &JWTService{
		jwtSigningKey: signingKey,
		defaultCookie: http.Cookie{
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Domain:   cookieDomain,
			Secure:   secureCookie,
		},
		jwtSessionLength: time.Duration(sessionLength),
		jwtSigningMethod: jwt.SigningMethodHS256,
	}
}

func (s *JWTService) GetUserIDFromRequest(r *http.Request) (uuid.UUID, error) {
	jwtCookie, err := r.Cookie(jwtCookieName)

	if err != nil {
		return uuid.Nil, err
	}

	userID, err := s.DecodeJWTToUser(jwtCookie.Value)

	if userID == uuid.Nil || err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (s *JWTService) CreateJWTTokenForUser(userID uuid.UUID) (string, error) {
	claims := CustomClaims{
		userID.String(),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * s.jwtSessionLength).Unix(),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSigningKey)
}

func (s *JWTService) DecodeJWTToUser(token string) (uuid.UUID, error) {
	decodeToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if !(s.jwtSigningMethod == token.Method) {
			return nil, errSigningMethodMismatch
		}

		return s.jwtSigningKey, nil
	})

	if decodedClaims, ok := decodeToken.Claims.(*CustomClaims); ok && decodeToken.Valid {
		if userID, err := uuid.Parse(decodedClaims.UserID); err != nil {
			return uuid.Nil, err
		} else {
			return userID, nil
		}
	}

	return uuid.Nil, err
}

func (s *JWTService) AuthCookie(token string) *http.Cookie {
	d := s.defaultCookie
	d.Name = jwtCookieName
	d.Value = token
	d.Path = "/"
	return &d
}

func (s *JWTService) ExpiredAuthCookie() *http.Cookie {
	d := s.defaultCookie
	d.Name = jwtCookieName
	d.Value = ""
	d.Path = "/"
	d.MaxAge = -1
	d.Expires = time.Date(1983, 7, 26, 20, 34, 58, 651387237, time.UTC)
	return &d
}
