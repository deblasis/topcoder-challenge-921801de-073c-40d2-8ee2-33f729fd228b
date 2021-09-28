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
package endpoints

import (
	"context"

	"github.com/sony/gobreaker"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/common/middlewares"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type EndpointSet struct {
	StatusEndpoint            endpoint.Endpoint
	GetUserByUsernameEndpoint endpoint.Endpoint
	GetUserByIdEndpoint       endpoint.Endpoint
	CreateUserEndpoint        endpoint.Endpoint
	logger                    log.Logger
}

func NewEndpointSet(s service.AuthDBService, logger log.Logger) EndpointSet {

	var getUserByUsernameEndpoint endpoint.Endpoint
	{
		getUserByUsernameEndpoint = MakeGetUserByUsernameEndpoint(s)
		getUserByUsernameEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getUserByUsernameEndpoint)
		getUserByUsernameEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "GetUserByUsername"))(getUserByUsernameEndpoint)
	}

	var getUserByIdEndpoint endpoint.Endpoint
	{
		getUserByIdEndpoint = MakeGetUserByIdEndpoint(s)
		getUserByIdEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getUserByIdEndpoint)
		getUserByIdEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "GetUserById"))(getUserByIdEndpoint)
	}

	var createUserEndpoint endpoint.Endpoint
	{
		createUserEndpoint = MakeCreateUserEndpoint(s)
		createUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(createUserEndpoint)
		createUserEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "CreateUser"))(createUserEndpoint)
	}

	return EndpointSet{
		StatusEndpoint:            healthcheck.MakeStatusEndpoint(logger),
		GetUserByUsernameEndpoint: getUserByUsernameEndpoint,
		GetUserByIdEndpoint:       getUserByIdEndpoint,
		CreateUserEndpoint:        createUserEndpoint,
		logger:                    logger,
	}
}

func MakeGetUserByIdEndpoint(s service.AuthDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.GetUserByIdRequest)
		resp, err := s.GetUserById(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeGetUserByUsernameEndpoint(s service.AuthDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.GetUserByUsernameRequest)
		resp, err := s.GetUserByUsername(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeCreateUserEndpoint(s service.AuthDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.CreateUserRequest)
		resp, err := s.CreateUser(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}
