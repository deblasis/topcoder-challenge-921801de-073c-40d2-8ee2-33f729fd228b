package endpoints

import (
	"context"

	"deblasis.net/space-traffic-control/common/errors"
	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/dtos"
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

// Signup(ctx context.Context, request dtos.SignupRequest) (dtos.SignupResponse, error)
func (s EndpointSet) Signup(ctx context.Context, request dtos.SignupRequest) (dtos.SignupResponse, error) {
	var ret dtos.SignupResponse
	resp, err := s.SignupEndpoint(ctx, dtos.SignupRequest{
		Username: request.Username,
		Password: request.Password,
		Role:     request.Role,
	})
	if err != nil {
		return ret, err
	}
	response := resp.(dtos.SignupResponse)
	return response, errors.Str2err(response.Err)
}

// Login(ctx context.Context, request dtos.LoginRequest) (dtos.LoginResponse, error)
func (s EndpointSet) Login(ctx context.Context, request dtos.LoginRequest) (dtos.LoginResponse, error) {
	var ret dtos.LoginResponse
	resp, err := s.LoginEndpoint(ctx, dtos.LoginRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return ret, err
	}
	response := resp.(dtos.LoginResponse)
	return response, errors.Str2err(response.Err)
}

var (
	_ endpoint.Failer = dtos.SignupResponse{}
	_ endpoint.Failer = dtos.LoginResponse{}
)
