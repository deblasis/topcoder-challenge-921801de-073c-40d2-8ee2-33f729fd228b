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
package acl

import (
	"fmt"

	"deblasis.net/space-traffic-control/common/auth"
	"deblasis.net/space-traffic-control/common/consts"
	"github.com/go-kit/kit/log"
	"github.com/golang-jwt/jwt"
)

func AclRules() map[string]auth.ACLRule {

	return map[string]auth.ACLRule{
		"/deblasis.v1.ShippingStationService/RequestLanding": func(req interface{}, logger log.Logger) auth.ACLDescriptor {

			return auth.NewMustCheckTokenDescriptor(func(t *jwt.Token) error {

				claims := t.Claims.(*auth.STCClaims)
				role := claims.Role

				if role != consts.ROLE_SHIP {
					return fmt.Errorf("unauthorized: you must be a %v in order to perform this operation", consts.ROLE_SHIP)
				}
				return nil

			})
		},
		"/deblasis.v1.ShippingStationService/Landing": func(req interface{}, logger log.Logger) auth.ACLDescriptor {

			return auth.NewMustCheckTokenDescriptor(func(t *jwt.Token) error {

				claims := t.Claims.(*auth.STCClaims)
				role := claims.Role

				if role != consts.ROLE_SHIP {
					return fmt.Errorf("unauthorized: you must be a %v in order to perform this operation", consts.ROLE_SHIP)
				}
				return nil

			})
		},
	}
}
