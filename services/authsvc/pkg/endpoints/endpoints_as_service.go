package endpoints

import (
	"context"

	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"github.com/go-kit/kit/endpoint"
)

// Signup(ctx context.Context, request pb.SignupRequest) (pb.SignupResponse, error)
func (s EndpointSet) Signup(ctx context.Context, request *pb.SignupRequest) (*pb.SignupResponse, error) {
	var ret *pb.SignupResponse
	resp, err := s.SignupEndpoint(ctx, request)
	if err != nil {
		return ret, err
	}
	response := resp.(*pb.SignupResponse)
	return response, nil
}

// Login(ctx context.Context, request pb.LoginRequest) (pb.LoginResponse, error)
func (s EndpointSet) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	var ret *pb.LoginResponse
	resp, err := s.LoginEndpoint(ctx, request)
	if err != nil {
		return ret, err
	}
	response := resp.(*pb.LoginResponse)
	return response, nil
	//return response, errors.Str2err(response.Error)
}

func (s EndpointSet) CheckToken(ctx context.Context, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	var ret *pb.CheckTokenResponse
	resp, err := s.CheckTokenEndpoint(ctx, request)
	if err != nil {
		return ret, err
	}
	response := resp.(*pb.CheckTokenResponse)
	return response, nil
	//return response, errors.Str2err(response.Error)
}

var (
	_ endpoint.Failer = pb.SignupResponse{}
	_ endpoint.Failer = pb.LoginResponse{}
	_ endpoint.Failer = pb.CheckTokenResponse{}
)
