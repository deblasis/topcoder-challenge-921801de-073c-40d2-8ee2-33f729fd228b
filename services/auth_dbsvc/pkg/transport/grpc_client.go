// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
package transport

import (
	"context"
	"strings"
	"time"

	"deblasis.net/space-traffic-control/common/errs"
	pb "deblasis.net/space-traffic-control/gen/proto/go/auth_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

func NewGRPCClient(conn *grpc.ClientConn, logger log.Logger) service.AuthDBService {

	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	var getUserByUsernameEndpoint endpoint.Endpoint
	{
		getUserByUsernameEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"GetUserByUsername",
			encodeGRPCGetUserByUsernameRequest,
			decodeGRPCGetUserByUsernameResponse,
			pb.GetUserByUsernameResponse{},
		).Endpoint()

		getUserByUsernameEndpoint = limiter(getUserByUsernameEndpoint)
		getUserByUsernameEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetUserByUsername",
			Timeout: 30 * time.Second,
		}))(getUserByUsernameEndpoint)
	}
	var createUserEndpoint endpoint.Endpoint
	{
		createUserEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"CreateUser",
			encodeGRPCCreateUserRequest,
			decodeGRPCCreateUserResponse,
			pb.CreateUserResponse{},
		).Endpoint()

		createUserEndpoint = limiter(createUserEndpoint)
		createUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "CreateUser",
			Timeout: 30 * time.Second,
		}))(createUserEndpoint)
	}

	return endpoints.EndpointSet{
		GetUserByUsernameEndpoint: getUserByUsernameEndpoint,
		CreateUserEndpoint:        createUserEndpoint,
	}

}

func encodeGRPCGetUserByUsernameRequest(_ context.Context, request interface{}) (interface{}, error) {
	//TODO converters
	req := request.(*dtos.GetUserByUsernameRequest)
	return &pb.GetUserByUsernameRequest{
		Username: req.Username,
	}, nil
}

func decodeGRPCGetUserByUsernameResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetUserByUsernameResponse)

	var user *model.User
	if response.User != nil {
		user = &model.User{
			Id:       response.User.Id,
			Username: response.User.Username,
			Password: response.User.Password,
			//Role:     strings.Title(strings.ToLower(strings.TrimLeft(response.User.Role.String(), "ROLE_"))),
			Role: response.User.Role,
		}
	}

	return &dtos.GetUserResponse{
		User:  user,
		Error: errs.FromProtoV1(response.Error),
	}, nil
}

func encodeGRPCCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*dtos.CreateUserRequest)

	return &pb.CreateUserRequest{
		User: &pb.User{
			Id:       req.Id,
			Username: req.Username,
			Password: req.Password,
			Role:     req.Role,
		},
	}, nil
}
func decodeGRPCCreateUserResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.CreateUserResponse)

	id, _ := uuid.Parse(response.Id)

	return &dtos.CreateUserResponse{
		Id:    &id,
		Error: errs.FromProtoV1(response.Error),
	}, nil
}
