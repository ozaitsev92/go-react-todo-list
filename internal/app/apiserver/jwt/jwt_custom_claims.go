package jwt

import (
	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}
