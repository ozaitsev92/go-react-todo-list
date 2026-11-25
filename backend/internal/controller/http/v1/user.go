package v1

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/jwt"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/model"
	"github.com/ozaitsev92/tododdd/internal/usecase"
)

type userRoutes struct {
	users  *usecase.UserUseCase
	tokens *jwt.Service
	log    *slog.Logger
}

type credentialsRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func newUserRoutes(g *gin.RouterGroup, users *usecase.UserUseCase, tokens *jwt.Service, auth gin.HandlerFunc, log *slog.Logger) {
	r := &userRoutes{users: users, tokens: tokens, log: log}

	g.POST("", r.register)
	g.POST("/login", r.login)
	g.POST("/logout", r.logout)
	g.GET("/current", auth, r.current)
}

func (r *userRoutes) register(c *gin.Context) {
	var req credentialsRequest
	if !bindJSON(c, &req) {
		return
	}

	u, err := r.users.RegisterNewUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		writeError(c, r.log, err)
		return
	}

	c.JSON(http.StatusCreated, model.ToResponseFromUser(u))
}

func (r *userRoutes) login(c *gin.Context) {
	var req credentialsRequest
	if !bindJSON(c, &req) {
		return
	}

	u, err := r.users.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		writeError(c, r.log, err)
		return
	}

	token, expires, err := r.tokens.Issue(u.ID)
	if err != nil {
		writeError(c, r.log, err)
		return
	}

	r.tokens.SetCookie(c.Writer, token, expires)
	c.JSON(http.StatusOK, model.ToResponseFromUser(u))
}

func (r *userRoutes) logout(c *gin.Context) {
	r.tokens.ClearCookie(c.Writer)
	c.Status(http.StatusNoContent)
}

func (r *userRoutes) current(c *gin.Context) {
	u, ok := currentUser(c)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, model.ToResponseFromUser(u))
}
