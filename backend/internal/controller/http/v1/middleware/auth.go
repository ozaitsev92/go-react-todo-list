// Package middleware contains the Gin middleware used by the v1 API.
package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
)

// TokenVerifier extracts the authenticated user ID from a request.
type TokenVerifier interface {
	UserIDFromRequest(r *http.Request) (uuid.UUID, error)
}

// UserGetter loads a user by ID.
type UserGetter interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (user.User, error)
}

type currentUserKey struct{}

// Authenticate answers 401 unless the session is valid and its user still
// exists, and stores that user in the request context (see CurrentUser).
func Authenticate(tokens TokenVerifier, users UserGetter, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := tokens.UserIDFromRequest(c.Request)
		if err != nil {
			abortJSON(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		u, err := users.GetUserByID(c.Request.Context(), userID)
		if errors.Is(err, user.ErrUserNotFound) {
			abortJSON(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		if err != nil {
			log.ErrorContext(c.Request.Context(), "auth: load user", "err", err, "user_id", userID)
			abortJSON(c, http.StatusInternalServerError, "internal server error")
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), currentUserKey{}, u))
		c.Next()
	}
}

// CurrentUser returns the user stored by Authenticate.
func CurrentUser(ctx context.Context) (user.User, bool) {
	u, ok := ctx.Value(currentUserKey{}).(user.User)
	return u, ok
}

func abortJSON(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, gin.H{"error": msg})
}
