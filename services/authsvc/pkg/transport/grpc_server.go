package transport

import (
	"context"

	"deblasis.net/space-traffic-control/common/healthcheck"
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	serviceStatus grpctransport.Handler
	signup        grpctransport.Handler
	login         grpctransport.Handler
}

func NewGRPCServer(e endpoints.EndpointSet, l log.Logger) pb.AuthServiceServer {
	return &grpcServer{
		serviceStatus: grpctransport.NewServer(
			e.StatusEndpoint,
			decodeGRPCServiceStatusRequest,
			encodeGRPCServiceStatusResponse,
		),
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
	}
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *pb.ServiceStatusRequest) (*pb.ServiceStatusResponse, error) {
	_, rep, err := g.serviceStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ServiceStatusResponse), nil
}
func decodeGRPCServiceStatusRequest(_ context.Context, _ interface{}) (interface{}, error) {
	var req healthcheck.ServiceStatusRequest
	return req, nil
}
func encodeGRPCServiceStatusResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	Response := grpcResponse.(healthcheck.ServiceStatusResponse)
	return &pb.ServiceStatusResponse{Code: Response.Code, Err: Response.Err}, nil
}

func (g *grpcServer) Signup(ctx context.Context, r *pb.SignupRequest) (*pb.SignupResponse, error) {
	_, rep, err := g.signup.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	res := rep.(pb.SignupResponse)
	return &res, nil
}
func (g *grpcServer) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, rep, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	res := rep.(pb.LoginResponse)
	return &res, nil
}

func decodeGRPCSignupRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SignupRequest)
	return req, nil
}
func encodeGRPCSignupResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(pb.SignupResponse)
	return response, nil
}

func decodeGRPCLoginRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LoginRequest)
	return req, nil
}

func encodeGRPCLoginResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(pb.LoginResponse)
	return response, nil
}
