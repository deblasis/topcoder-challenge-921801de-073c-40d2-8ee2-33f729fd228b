package auth

import (
	"time"

	"deblasis.net/space-traffic-control/common/config"
	"github.com/golang-jwt/jwt"
)

type STCClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func NewJWTToken(config config.JWTConfig, username, role, issuer string) (token string, expiresAt int64, err error) {
	now := time.Now().Unix()
	expiresAt = now + int64(config.TokenDuration)
	claims := STCClaims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			Audience:  "deblasis.SpaceTrafficControl",
			ExpiresAt: expiresAt,
			IssuedAt:  now,
			Issuer:    issuer,
		},
	}

	token, err = generateToken(config, claims)
	return token, expiresAt, err
}

func generateToken(config config.JWTConfig, claims STCClaims) (string, error) {

	//Simple implementation here, ideally we would pick this up from Vault and/or we'll use RSA or something more secure
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Secret))

	return tokenString, err
}
