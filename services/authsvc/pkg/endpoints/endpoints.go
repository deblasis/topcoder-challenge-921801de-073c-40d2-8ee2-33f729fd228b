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

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/middlewares"
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
)

type EndpointSet struct {
	SignupEndpoint     endpoint.Endpoint
	LoginEndpoint      endpoint.Endpoint
	CheckTokenEndpoint endpoint.Endpoint

	logger log.Logger
}

func NewEndpointSet(s service.AuthService, logger log.Logger) EndpointSet {

	var signupEndpoint endpoint.Endpoint
	{
		signupEndpoint = MakeSignupEndpoint(s)
		signupEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(signupEndpoint)
		signupEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "Signup"))(signupEndpoint)
	}

	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = MakeLoginEndpoint(s)
		loginEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(loginEndpoint)
		loginEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "Login"))(loginEndpoint)
	}

	var checkTokenEndpoint endpoint.Endpoint
	{
		checkTokenEndpoint = MakeCheckTokenEndpoint(s)
		checkTokenEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(checkTokenEndpoint)
		checkTokenEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "CheckToken"))(checkTokenEndpoint)
	}

	return EndpointSet{
		SignupEndpoint:     signupEndpoint,
		LoginEndpoint:      loginEndpoint,
		CheckTokenEndpoint: checkTokenEndpoint,
		logger:             logger,
	}
}

func MakeSignupEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.SignupRequest)
		resp, err := s.Signup(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeLoginEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.LoginRequest)
		resp, err := s.Login(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeCheckTokenEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.CheckTokenRequest)
		resp, err := s.CheckToken(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}
