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
	"github.com/go-kit/log"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/codes"
)

func AclRules() map[string]auth.ACLRule {

	return map[string]auth.ACLRule{
		"/deblasis.v1.CentralCommandService/RegisterShip": func(req interface{}, logger log.Logger) auth.ACLDescriptor {

			return auth.NewMustCheckTokenDescriptor(func(t *jwt.Token) error {

				claims := t.Claims.(*auth.STCClaims)
				role := claims.Role

				if role != consts.ROLE_SHIP {
					return fmt.Errorf("unauthorized: you must be a %v in order to perform this operation", consts.ROLE_SHIP)
				}
				return nil

			})
		},
		"/deblasis.v1.CentralCommandService/RegisterStation": func(req interface{}, logger log.Logger) auth.ACLDescriptor {

			return auth.NewMustCheckTokenDescriptor(func(t *jwt.Token) error {

				claims := t.Claims.(*auth.STCClaims)
				role := claims.Role

				if role != consts.ROLE_STATION {
					return fmt.Errorf("unauthorized: you must be a %v in order to perform this operation", consts.ROLE_SHIP)
				}
				return nil

			})
		},
		"/deblasis.v1.CentralCommandService/GetAllStations": func(req interface{}, logger log.Logger) auth.ACLDescriptor {

			return auth.NewMustCheckTokenDescriptorWithCustomStatusCode(func(t *jwt.Token) error {

				claims := t.Claims.(*auth.STCClaims)
				role := claims.Role

				if role != consts.ROLE_COMMAND && role != consts.ROLE_SHIP {
					return fmt.Errorf("unauthorized: you must be a %v or %v in order to perform this operation", consts.ROLE_COMMAND, consts.ROLE_SHIP)
				}
				return nil
			}, codes.InvalidArgument)
		},
		"/deblasis.v1.CentralCommandService/GetAllShips": func(req interface{}, logger log.Logger) auth.ACLDescriptor {

			return auth.NewMustCheckTokenDescriptor(func(t *jwt.Token) error {

				claims := t.Claims.(*auth.STCClaims)
				role := claims.Role

				if role != consts.ROLE_COMMAND {
					return fmt.Errorf("unauthorized: you must be a %v in order to perform this operation", consts.ROLE_COMMAND)
				}
				return nil
			})
		},
	}
}

// func(req interface{}) auth.TokenValidator {
// 	request := req.(*pb.SignupRequest)

// 	if request.Role == consts.ROLE_SHIP {
// 		return struct{
// 			allGood
// 		}
// 	}

// 	return func(t *jwt.Token) error {

// 		claims := t.Claims.(*auth.STCClaims)
// 		role := claims.Role

// 		if request.Role == consts.ROLE_SHIP || role == consts.ROLE_COMMAND {
// 			return nil
// 		}
// 		return fmt.Errorf("unauthorized: you must be a member of %v to perform this operation", consts.ROLE_COMMAND)
// 	}
// },
