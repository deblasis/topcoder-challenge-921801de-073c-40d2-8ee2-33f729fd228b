//
// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
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
