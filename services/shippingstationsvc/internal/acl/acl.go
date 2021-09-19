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

			return auth.NewMustCheckTokenDescriptor(func(t *jwt.Token) error {

				// claims := t.Claims.(jwt.MapClaims) //auth.STCClaims
				// role := claims["role"]
				claims := t.Claims.(*auth.STCClaims)
				role := claims.Role

				if role != consts.ROLE_COMMAND && role != consts.ROLE_SHIP {
					return fmt.Errorf("unauthorized: you must be a %v or %v in order to perform this operation", consts.ROLE_COMMAND, consts.ROLE_SHIP)
				}
				return nil
			})
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
