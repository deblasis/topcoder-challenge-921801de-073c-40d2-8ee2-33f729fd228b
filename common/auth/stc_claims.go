package auth

import "github.com/golang-jwt/jwt"

type STCClaims struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
