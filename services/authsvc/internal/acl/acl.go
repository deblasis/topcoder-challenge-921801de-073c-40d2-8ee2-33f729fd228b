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
