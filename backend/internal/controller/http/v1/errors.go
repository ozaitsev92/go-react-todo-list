package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/middleware"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	"github.com/ozaitsev92/tododdd/internal/usecase"
)

// errorResponse is the API's only error body.
type errorResponse struct {
	Error string `json:"error"`
}

func abortError(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, errorResponse{Error: msg})
}

// writeError maps errors to responses; unknown ones are logged and hidden behind a 500.
func writeError(c *gin.Context, log *slog.Logger, err error) {
	switch {
	// Another user's task looks exactly like a missing one.
	case errors.Is(err, task.ErrTaskNotFound), errors.Is(err, usecase.ErrUnauthorizedAction):
		abortError(c, http.StatusNotFound, "task not found")
	case errors.Is(err, task.ErrInvalidText),
		errors.Is(err, user.ErrInvalidEmail),
		errors.Is(err, user.ErrWeakPassword):
		abortError(c, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, user.ErrEmailTaken):
		abortError(c, http.StatusConflict, err.Error())
	case errors.Is(err, usecase.ErrInvalidCredentials):
		abortError(c, http.StatusUnauthorized, err.Error())
	default:
		log.ErrorContext(c.Request.Context(), "request failed",
			"err", err, "method", c.Request.Method, "route", c.FullPath())
		abortError(c, http.StatusInternalServerError, "internal server error")
	}
}

// bindJSON decodes the body into dst; on failure it writes 413 or 400 and returns false.
func bindJSON(c *gin.Context, dst any) bool {
	err := c.ShouldBindJSON(dst)
	if err == nil {
		return true
	}

	if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
		abortError(c, http.StatusRequestEntityTooLarge, "request body too large")
		return false
	}

	abortError(c, http.StatusBadRequest, "invalid request body")

	return false
}

// pathID parses the :id parameter, writing 400 if it is not a UUID.
func pathID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		abortError(c, http.StatusBadRequest, "invalid task id")
		return uuid.Nil, false
	}

	return id, true
}

// currentUser returns the user set by middleware.Authenticate, writing 401 if absent.
func currentUser(c *gin.Context) (user.User, bool) {
	u, ok := middleware.CurrentUser(c.Request.Context())
	if !ok {
		abortError(c, http.StatusUnauthorized, "unauthorized")
	}

	return u, ok
}
