package auth

import (
	"context"

	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
)

type TokenVerifier interface {
	CheckToken(ctx context.Context, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error)
}
