package jwt

import (
	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}
