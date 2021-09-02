package endpoints

import (
	"context"

	"deblasis.net/space-traffic-control/services/auth_dbsvc/api/model"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type EndpointSet struct {
	StatusEndpoint            endpoint.Endpoint
	GetUserByUsernameEndpoint endpoint.Endpoint
	CreateUserEndpoint        endpoint.Endpoint
	logger                    log.Logger
}

func NewEndpointSet(s service.UserManager, logger log.Logger) EndpointSet {
	return EndpointSet{
		StatusEndpoint:            makeStatusEndpoint(s),
		GetUserByUsernameEndpoint: makeGetUserByUsernameEndpoint(s),
		CreateUserEndpoint:        makeCreateUserEndpointEndpoint(s),
		logger:                    logger,
	}
}

func makeStatusEndpoint(s service.UserManager) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(model.ServiceStatusRequest)
		code, err := s.ServiceStatus(ctx)
		if err != nil {
			return model.ServiceStatusReply{Code: code, Err: err.Error()}, nil
		}
		return model.ServiceStatusReply{Code: code, Err: ""}, nil
	}
}

func makeGetUserByUsernameEndpoint(s service.UserManager) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.GetUserByUsernameRequest)
		user, err := s.GetUserByUsername(ctx, req.Username)
		if err != nil {
			return model.GetUserByUsernameReply{
				User: user,
				Err:  err.Error(),
			}, nil
		}
		return model.GetUserByUsernameReply{
			User: user,
			Err:  "",
		}, nil
	}
}

func makeCreateUserEndpointEndpoint(s service.UserManager) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.CreateUserRequest)
		id, err := s.CreateUser(ctx, &req.User)
		if err != nil {
			return model.CreateUserReply{
				Id:  -1,
				Err: err.Error(),
			}, err
		}
		return model.CreateUserReply{
			Id:  id,
			Err: "",
		}, nil
	}
}
