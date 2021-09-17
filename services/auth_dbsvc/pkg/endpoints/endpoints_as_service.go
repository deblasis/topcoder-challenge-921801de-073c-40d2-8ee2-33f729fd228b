package endpoints

import (
	"context"

	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	"github.com/go-kit/kit/endpoint"
)

// GetUserByUsername(ctx context.Context, username string) (model.User, error)
func (s EndpointSet) GetUserByUsername(ctx context.Context, request *dtos.GetUserByUsernameRequest) (*dtos.GetUserResponse, error) {
	resp, err := s.GetUserByUsernameEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.GetUserResponse)
	return response, nil
}

// GetUserById(ctx context.Context, id string) (model.User, error)
func (s EndpointSet) GetUserById(ctx context.Context, request *dtos.GetUserByIdRequest) (*dtos.GetUserResponse, error) {
	resp, err := s.GetUserByIdEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.GetUserResponse)
	return response, nil
}

// CreateUser(ctx context.Context, user *model.User) (int64, error)
func (s EndpointSet) CreateUser(ctx context.Context, request *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error) {
	resp, err := s.CreateUserEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.CreateUserResponse)
	return response, nil
}

var (
	_ endpoint.Failer = dtos.GetUserResponse{}
	_ endpoint.Failer = dtos.CreateUserResponse{}
)
