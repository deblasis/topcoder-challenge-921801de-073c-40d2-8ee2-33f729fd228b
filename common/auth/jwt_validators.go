// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>  
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE. 
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
