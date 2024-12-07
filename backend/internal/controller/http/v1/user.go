package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/jwt"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/middleware"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/model"
	"github.com/ozaitsev92/tododdd/internal/usecase"
	"github.com/ozaitsev92/tododdd/pkg/logger"
)

type userRoutes struct {
	l          logger.Interface
	jwtService *jwt.JWTService
	u          *usecase.UserUseCase
}

// todo: refactor. too many params
func newUserRoutes(handler *gin.RouterGroup, l logger.Interface, jwtService *jwt.JWTService, u *usecase.UserUseCase) {
	r := &userRoutes{l, jwtService, u}

	h := handler.Group("/users")
	{
		jwtMiddleware := middleware.JwtMiddleware(u, jwtService)

		h.POST("", r.createUser)
		h.POST("/login", r.loginUser)
		h.POST("/logout", r.logoutUser)
		h.GET("/current", jwtMiddleware, r.currentUser)
	}
}

func (r *userRoutes) createUser(c *gin.Context) {
	type createUserRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var request createUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - createUser")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})

		return
	}

	u, err := r.u.RegisterNewUser(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		r.l.Error(err, "http - v1 - createUser")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})

		return
	}

	c.JSON(http.StatusCreated, model.ToResponseFromUser(u))
}

func (r *userRoutes) loginUser(c *gin.Context) {
	type loginUserRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var request loginUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - loginUser")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})

		return
	}

	u, err := r.u.GetUserByEmail(c.Request.Context(), request.Email)
	if err != nil {
		r.l.Error(err, "http - v1 - loginUser")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Invalid login and/or password"})

		return
	}

	if !u.ComparePassword(request.Password) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Invalid login and/or password"})

		return
	}

	token, err := r.jwtService.CreateJWTTokenForUser(u.ID)
	if err != nil {
		r.l.Error(err, "http - v1 - loginUser")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "Internal server error"})

		return
	}
	// todo: add a refresh token
	// https://medium.com/novai-go-programming-101/building-a-jwt-authentication-system-with-refresh-tokens-in-go-adce3b30c1ac
	cookie := r.jwtService.AuthCookie(token)

	c.SetCookie(
		cookie.Name,
		cookie.Value,
		cookie.MaxAge,
		cookie.Path,
		cookie.Domain,
		cookie.Secure,
		cookie.HttpOnly,
	)
	c.JSON(http.StatusOK, gin.H{})
}

func (r *userRoutes) logoutUser(c *gin.Context) {
	cookie := r.jwtService.ExpiredAuthCookie()

	c.SetCookie(
		cookie.Name,
		cookie.Value,
		cookie.MaxAge,
		cookie.Path,
		cookie.Domain,
		cookie.Secure,
		cookie.HttpOnly,
	)
	c.JSON(http.StatusOK, gin.H{})
}

func (r *userRoutes) currentUser(c *gin.Context) {
	id := c.GetString("userID")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})

		return
	}

	u, err := r.u.GetUserByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})

		return
	}

	c.JSON(http.StatusOK, model.ToResponseFromUser(u))
}
