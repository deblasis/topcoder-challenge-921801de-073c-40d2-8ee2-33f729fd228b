package auth

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

type TokenValidator func(*jwt.Token) error

func defaultValidator(token *jwt.Token) error {
	//here some validation we ALWAYS perform on a valid token

	claims := token.Claims.(STCClaims)

	if claims.Role == "" {
		return errors.New("role cannot be empty")
	}

	return nil
}

func AllowRoles(roles ...string) TokenValidator {

	return func(token *jwt.Token) error {
		claims := token.Claims.(STCClaims)

		for _, allowedRole := range roles {
			if allowedRole == claims.Role {
				return nil
			}
		}

		return errors.New("role is not allowed to pass")
	}
}
