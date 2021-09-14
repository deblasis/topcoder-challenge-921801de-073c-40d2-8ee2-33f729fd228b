package transport

import (
	"context"

	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	signup     grpctransport.Handler
	login      grpctransport.Handler
	checkToken grpctransport.Handler
}

func NewGRPCServer(l log.Logger, e endpoints.EndpointSet) pb.AuthServiceServer {
	return &grpcServer{
		signup: grpctransport.NewServer(
			e.SignupEndpoint,
			decodeGRPCSignupRequest,
			encodeGRPCSignupResponse,
		),
		login: grpctransport.NewServer(
			e.LoginEndpoint,
			decodeGRPCLoginRequest,
			encodeGRPCLoginResponse,
		),
		checkToken: grpctransport.NewServer(
			e.CheckTokenEndpoint,
			decodeGRPCCheckTokenRequest,
			encodeGRPCCheckTokenResponse,
		),
	}
}

func (g *grpcServer) Signup(ctx context.Context, r *pb.SignupRequest) (*pb.SignupResponse, error) {
	_, rep, err := g.signup.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SignupResponse), nil
}

func decodeGRPCSignupRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.SignupRequest), nil
}
func encodeGRPCSignupResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	return grpcResponse.(*pb.SignupResponse), nil
}

func (g *grpcServer) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, rep, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.LoginResponse), nil
}
func decodeGRPCLoginRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.LoginRequest), nil
}

func encodeGRPCLoginResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	return grpcResponse.(*pb.LoginResponse), nil
}

func (g *grpcServer) CheckToken(ctx context.Context, r *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	_, rep, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CheckTokenResponse), nil
}
func decodeGRPCCheckTokenRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.CheckTokenRequest), nil
}

func encodeGRPCCheckTokenResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	return grpcResponse.(*pb.CheckTokenResponse), nil
}
