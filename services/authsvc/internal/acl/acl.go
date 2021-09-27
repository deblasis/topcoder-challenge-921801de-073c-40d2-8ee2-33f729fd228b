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
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"github.com/go-kit/kit/log"
	"github.com/golang-jwt/jwt"
)

func AclRules() map[string]auth.ACLRule {

	return map[string]auth.ACLRule{
		"/deblasis.v1.AuthService/Signup": func(req interface{}, logger log.Logger) auth.ACLDescriptor {

			request := req.(*pb.SignupRequest)
			if request.Role == consts.ROLE_SHIP {
				return auth.NewAllGoodAclDescriptor()
			}
			return auth.NewMustCheckTokenDescriptor(func(t *jwt.Token) error {

				// claims := t.Claims.(jwt.MapClaims) //auth.STCClaims
				// role := claims["role"]
				claims := t.Claims.(*auth.STCClaims)
				role := claims.Role

				if role == consts.ROLE_COMMAND {
					return nil
				}
				return fmt.Errorf("unauthorized: you must be a member of %v to perform this operation", consts.ROLE_COMMAND)

			})
		},
	}
}
