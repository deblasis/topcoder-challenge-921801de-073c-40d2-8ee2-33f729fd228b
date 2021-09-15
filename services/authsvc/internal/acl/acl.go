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
		"/deblasis.v1.CentralCommandService/RegisterShip": func(req interface{}, logger log.Logger) auth.ACLDescriptor {

			return auth.NewMustCheckTokenDescriptor(func(t *jwt.Token) error {

				claims := t.Claims.(jwt.MapClaims) //auth.STCClaims
				role := claims["role"]

				if role != consts.ROLE_SHIP {
					return fmt.Errorf("unauthorized: you must be a %v in order to perform this operation", consts.ROLE_SHIP)
				}

				//check if already registered

				return nil
			})
		},
	}
}
