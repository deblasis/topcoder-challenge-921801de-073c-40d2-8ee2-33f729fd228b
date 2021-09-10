package converters

import (
	"strings"

	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
)

func ProtoToDTORole(src pb.SignupRequest_Role) string {
	return strings.Title(strings.TrimLeft(strings.ToLower(src.String()), "role_"))
}
