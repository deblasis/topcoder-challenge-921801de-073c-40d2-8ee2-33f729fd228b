package converters

// import (
// 	"strings"

// 	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
// )

// func ProtoToDTORole(src pb.SignupRequest_Role) string {
// 	return strings.Title(strings.TrimLeft(strings.ToLower(src.String()), "role_"))
// }

import (
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	m "gopkg.in/jeevatkm/go-model.v1"
)

func SignupRequestToDBDto(src *pb.SignupRequest) *dtos.CreateUserRequest {
	ret := &dtos.CreateUserRequest{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}
