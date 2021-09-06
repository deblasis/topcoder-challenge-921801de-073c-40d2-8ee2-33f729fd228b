package endpoints

import (
	"context"

	"deblasis.net/space-traffic-control/common/errors"
	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	"github.com/go-kit/kit/endpoint"
)

// ServiceStatus(ctx context.Context) (int64, error)
func (s EndpointSet) ServiceStatus(ctx context.Context) (int64, error) {
	resp, err := s.StatusEndpoint(ctx, healthcheck.ServiceStatusRequest{})
	if err != nil {
		return 0, err
	}
	response := resp.(healthcheck.ServiceStatusResponse)
	return response.Code, errors.Str2err(response.Err)
}

// GetUserByUsername(ctx context.Context, username string) (model.User, error)
func (s EndpointSet) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	resp, err := s.GetUserByUsernameEndpoint(ctx, dtos.GetUserByUsernameRequest{Username: username})
	if err != nil {
		return user, err
	}
	response := resp.(dtos.GetUserByUsernameResponse)
	return response.User, errors.Str2err(response.Err)
}

// CreateUser(ctx context.Context, user *model.User) (int64, error)
func (s EndpointSet) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	resp, err := s.CreateUserEndpoint(ctx, dtos.CreateUserRequest{
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	})
	if err != nil {
		return -1, err
	}
	response := resp.(dtos.CreateUserResponse)
	return response.Id, errors.Str2err(response.Err)
}

var (
	_ endpoint.Failer = dtos.GetUserByUsernameResponse{}
	_ endpoint.Failer = dtos.CreateUserResponse{}
)

// type ServiceStatusRequest struct{}
// type ServiceStatusResponse struct {
// 	Code int64 `json:"code"`
// 	Err  error `json:"-"`
// }

// type GetUserByUsernameRequest struct {
// 	Username string `json:"username" validate:"required,notblank"`
// }
// type GetUserByUsernameResponse struct {
// 	User model.User `json:"user"`
// 	Err  error      `json:"-"`
// }

// type CreateUserRequest model.User
// type CreateUserResponse struct {
// 	Id  int64 `json:"id"`
// 	Err error `json:"-"`
// }

//func (r ServiceStatusResponse) Failed() error { return r.Err }
//func (r GetUserByUsernameResponse) Failed() error { return r.Err }
//func (r CreateUserResponse) Failed() error { return r.Err }
