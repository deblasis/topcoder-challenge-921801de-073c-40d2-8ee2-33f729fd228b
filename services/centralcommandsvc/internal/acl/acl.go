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
package acl

import (
	"fmt"

	"deblasis.net/space-traffic-control/common/auth"
	"deblasis.net/space-traffic-control/common/consts"
	"github.com/go-kit/kit/log"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/codes"
)

func AclRules() map[string]auth.ACLRule {

	return map[string]auth.ACLRule{
		"/deblasis.v1.CentralCommandService/RegisterShip": func(req interface{}, logger log.Logger) auth.ACLDescriptor {

			return auth.NewMustCheckTokenDescriptor(func(t *jwt.Token) error {

				// claims := t.Claims.(jwt.MapClaims) //auth.STCClaims
				// role := claims["role"]
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

				// claims := t.Claims.(jwt.MapClaims) //auth.STCClaims
				// role := claims["role"]
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

				// claims := t.Claims.(jwt.MapClaims) //auth.STCClaims
				// role := claims["role"]
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

				// claims := t.Claims.(jwt.MapClaims) //auth.STCClaims
				// role := claims["role"]
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
