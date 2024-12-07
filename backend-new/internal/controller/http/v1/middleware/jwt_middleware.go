package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/jwt"
	"github.com/ozaitsev92/tododdd/internal/usecase"
)

func JwtMiddleware(u *usecase.UserUseCase, jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := jwtService.GetUserIDFromRequest(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})

			return
		}

		u, err := u.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})

			return
		}

		c.Set("user", u)

		c.Next()
	}
}
