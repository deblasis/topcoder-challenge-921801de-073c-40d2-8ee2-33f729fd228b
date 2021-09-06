package transport

import (
	"context"

	"strings"

	"deblasis.net/space-traffic-control/common/healthcheck"
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/dtos"
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
	return rep.(*pb.SignupResponse), nil
}
func (g *grpcServer) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, rep, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.LoginResponse), nil
}

func decodeGRPCSignupRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SignupRequest)
	return dtos.SignupRequest{
		Username: req.Username,
		Password: req.Password,
		//TODO centralize
		Role: strings.Title(strings.ToLower(strings.TrimLeft(req.Role.String(), "ROLE_"))),
	}, nil
}
func encodeGRPCSignupResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(dtos.SignupResponse)
	return &pb.SignupResponse{
		Token: &pb.Token{
			Token:     response.Token.Token,
			ExpiresAt: response.Token.ExpiresAt,
		},
		Error: response.Err,
	}, nil
}

func decodeGRPCLoginRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LoginRequest)
	return dtos.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}, nil
}

func encodeGRPCLoginResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(dtos.LoginResponse)
	return &pb.LoginResponse{
		Token: &pb.Token{
			Token:     response.Token.Token,
			ExpiresAt: response.Token.ExpiresAt,
		},
		Error: response.Err,
	}, nil
}
