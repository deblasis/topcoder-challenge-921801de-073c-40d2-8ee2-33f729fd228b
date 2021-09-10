package endpoints

import (
	"context"

	"deblasis.net/space-traffic-control/common/errors"
	"deblasis.net/space-traffic-control/common/healthcheck"
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"github.com/go-kit/kit/endpoint"
)

//ServiceStatus(ctx context.Context) (int64, error)
func (s EndpointSet) ServiceStatus(ctx context.Context) (int64, error) {
	resp, err := s.StatusEndpoint(ctx, healthcheck.ServiceStatusRequest{})
	if err != nil {
		return 0, err
	}
	response := resp.(healthcheck.ServiceStatusResponse)
	return response.Code, errors.Str2err(response.Err)
}

// Signup(ctx context.Context, request pb.SignupRequest) (pb.SignupResponse, error)
func (s EndpointSet) Signup(ctx context.Context, request pb.SignupRequest) (pb.SignupResponse, error) {
	var ret pb.SignupResponse
	resp, err := s.SignupEndpoint(ctx, pb.SignupRequest{
		Username: request.Username,
		Password: request.Password,
		Role:     request.Role,
	})
	if err != nil {
		return ret, err
	}
	response := resp.(pb.SignupResponse)
	return response, errors.Str2err(response.Error)
}

// Login(ctx context.Context, request pb.LoginRequest) (pb.LoginResponse, error)
func (s EndpointSet) Login(ctx context.Context, request pb.LoginRequest) (pb.LoginResponse, error) {
	var ret pb.LoginResponse
	resp, err := s.LoginEndpoint(ctx, pb.LoginRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return ret, err
	}
	response := resp.(pb.LoginResponse)
	return response, errors.Str2err(response.Error)
}

var (
	_ endpoint.Failer = pb.SignupResponse{}
	_ endpoint.Failer = pb.LoginResponse{}
)
